package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	// TokenLifespan until 14 days
	TokenLifespan = time.Hour * 24 * 14
	// KeyAuthUserID to use in context
	KeyAuthUserID key = "auth_user_id"
)

var (
	// ErrUnauthenticated used when ther is no authenticated user in context
	ErrUnauthenticated = errors.New("unauthenticated")
)

type key string

//LoginOutput response
type LoginOutput struct {
	Token     string    `json:"token,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	User      User      `json:"user,omitempty"`
}

//AuthUser ID from Token
func (s *Service) AuthUserID(token string) (int64, error) {
	str, err := s.codec.DecodeToString(token)
	if err != nil {
		return 0, fmt.Errorf("could not decode token: %v", err)
	}

	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("could not parse auth user id from token: %v", err)
	}
	return i, nil
}

//login insecurely
func (s *Service) Login(ctx context.Context, email string) (LoginOutput, error) {
	//계정 보안용
	//실제 프로젝트를 위해서 필요한 부분, 토이프로젝트 시 생략 가능
	var out LoginOutput

	email = strings.TrimSpace(email)
	if !rxEmail.MatchString(email) {
		return out, ErrInvalidEmail
	}

	query := "SELECT id, username FROM users WHERE email = $1"
	err := s.db.QueryRowContext(ctx, query, email).Scan(&out.User.ID, &out.User.UserName)

	if err == sql.ErrNoRows {
		return out, ErrUserNotFound
	}

	if err != nil {
		return out, fmt.Errorf("could not query select user: %v\n", err)
	}
	//유저 아이디 토큰화
	out.Token, err = s.codec.EncodeToString(strconv.FormatInt(int64(out.User.ID), 10))
	if err != nil {
		return out, fmt.Errorf("could not query select user: %v\n", err)
	}

	out.ExpiresAt = time.Now().Add(TokenLifespan)

	return out, nil
}

//AuthUser From context
func (s *Service) AuthUser(ctx context.Context) (User, error) {
	var u User
	uid, ok := ctx.Value(KeyAuthUserID).(int64)
	if !ok {
		return u, ErrUnauthenticated
	}

	//Database query with the user ID in the context
	query := "SELECT username FROM users WHERE id = $1"
	err := s.db.QueryRowContext(ctx, query, uid).Scan(&u.UserName)
	if err == sql.ErrNoRows {
		return u, ErrUserNotFound
	}

	if err != nil {
		return u, fmt.Errorf("could not query select auth user: %v", err)
	}

	u.ID = uid
	return u, nil
}
