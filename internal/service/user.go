package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	//이메일 및 이름 확인(정규식 사용)
	rxEmail    = regexp.MustCompile("^[^\\s@]+@[^\\s@]+\\.[^\\s@]+$")
	rxUsername = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_-]{0,17}$")
	// ErrUserNotFound used when the user wasn't found on the db.
	ErrUserNotFound = errors.New("user not found")
	// ErrInvalidEmail used when the mail. is not valid
	ErrInvalidEmail = errors.New("invalid email")
	// ErrInvalidUsername used when the name. is not valid
	ErrInvalidUsername = errors.New("invalid name")
	// ErrEmailTaken used when there is already an user registered with that email.
	ErrEmailTaken = errors.New("email taken")
	// ErrUsernameTaken used when there is already an user registered with that username.
	ErrUsernameTaken = errors.New("username taken")
)

//User Model
type User struct {
	ID       int64  `json:"id,omitempty"`
	UserName string `json:"user_name,omitempty"`
}

// CreateUser inserts a user int the database.
func (s *Service) CreateUser(ctx context.Context, email, username string) error {
	email = strings.TrimSpace(email)
	if !rxEmail.MatchString(email) {
		return ErrInvalidEmail
	}

	username = strings.TrimSpace(username)
	if !rxUsername.MatchString(username) {
		return ErrInvalidUsername
	}

	query := "INSERT INTO users (email, username) VALUES ($1, $2)"
	_, err := s.db.ExecContext(ctx, query, email, username)
	unique := isUniqueViolation(err)

	//동일한 데이터를 입력했을때
	if unique && strings.Contains(err.Error(), "email") {
		return ErrEmailTaken
	}

	if unique && strings.Contains(err.Error(), "username") {
		return ErrUsernameTaken
	}

	//그 외 오류
	if err != nil {
		return fmt.Errorf("could not insert user: %v", err)
	}

	return nil
}
