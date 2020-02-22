package service

import (
	"context"
	"database/sql"
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
	// ErrForbiddenFollow is used when you try to following yourself
	ErrForbiddenFollow = errors.New("cannot follow yourself")
)

//User Model
type User struct {
	ID       int64  `json:"id,omitempty"`
	UserName string `json:"user_name,omitempty"`
}

//팔로워 카운트
//ToggleFollowOutput response
type ToggleFollowOutput struct {
	Following      bool `json:"following,omitempty"`
	FollowersCount int  `json:"followers_count,omitempty"`
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

// 팔로워 이름을 받는 기능
// ToggleFollow between Two users.
func (s *Service) ToggleFollow(ctx context.Context, username string) (ToggleFollowOutput, error) {
	var out ToggleFollowOutput
	// 유저 아이디 체크
	// check for the auth user in the context
	followerID, ok := ctx.Value(KeyAuthUserID).(int64)
	if !ok {
		return out, ErrUnauthenticated
	}

	// 유저 이름 확인
	// validate username
	username = strings.TrimSpace(username)
	if !rxUsername.MatchString(username) {
		return out, ErrInvalidUsername
	}

	// 쿼리를 트랜잭션으로 보내기
	// for queries to use transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return out, fmt.Errorf("could not begin tx: %v", err)
	}

	defer tx.Rollback()

	// 유저 이름으로부터 유저 아이디 받아오기
	// Get the actual user ID from the username
	var followeeID int64
	query := "SELECT id FROM users WHERE username = $1"
	tx.QueryRowContext(ctx, query, username).Scan(&followeeID)
	if err == sql.ErrNoRows {
		return out, ErrUserNotFound
	}

	if err != nil {
		return out, fmt.Errorf("could not query select user id from followee username: %v", err)
	}

	if followeeID == followerID {
		return out, ErrForbiddenFollow
	}

	// check overlap ID
	// 팔로워 중복 확인
	query = "SELECT EXISTS (SELECT 1 FROM follows WHERE follower_id = $1 AND followee_id = $2)"
	if err = tx.QueryRowContext(ctx, query, followerID, followeeID).Scan(&out.Following); err != nil {
		return out, fmt.Errorf("could not query select existance of follow: %v", err)
	}

	// 팔로우를 두 번 눌렀을 때의 액션
	// Action when pressed follow button twice

	// 팔로우 취소와 팔로워 수 감소
	// Cancle follow and decrease follow
	if out.Following {
		//취소
		//cancle
		query = "DELETE FROM follows WHERE follower_id = $1 AND followee_id = $2"
		if _, err = tx.ExecContext(ctx, query, followerID, followeeID); err != nil {
			return out, fmt.Errorf("could not delete follow : %v", err)
		}
		//팔로위 감소
		//followee
		query = "UPDATE users SET followees_count = followees_count - 1 WHERE id = $1"
		if _, err = tx.ExecContext(ctx, query, followerID); err != nil {
			return out, fmt.Errorf("could not update follower followees count (-): %v", err)
		}
		//팔로워 감소
		//follower
		query = "UPDATE users SET followers_count = followers_count - 1 WHERE id = $1 RETURNING followers_count"
		if err = tx.QueryRowContext(ctx, query, followeeID).Scan(&out.FollowersCount); err != nil {
			return out, fmt.Errorf("could not update followee followers count (-): %v", err)
		}
	} else {
		// 팔로우와 팔로워 증가
		// increment follow and follower
		query = "INSERT INTO follows (follower_id, followee_id) VALUES ($1, $2)"
		if _, err = tx.ExecContext(ctx, query, followerID, followeeID); err != nil {
			return out, fmt.Errorf("could not insert follow : %v", err)
		}
		// 팔로위 증가
		// followee increment
		query = "UPDATE users SET followees_count = followees_count + 1 WHERE id = $1"
		if _, err = tx.ExecContext(ctx, query, followerID); err != nil {
			return out, fmt.Errorf("could not update follower followees count (+): %v", err)
		}
		// 팔로워 증가
		// follower increment
		query = "UPDATE users SET followers_count = followers_count + 1 WHERE id = $1 RETURNING followers_count"
		if err = tx.QueryRowContext(ctx, query, followeeID).Scan(&out.FollowersCount); err != nil {
			return out, fmt.Errorf("could not update followee followers count (+): %v", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return out, fmt.Errorf("could not commit toggle follow : %v", err)
	}

	out.Following = !out.Following
	// 팔로잉을 끊었을 시 더이상 알림 받지 않음
	// Dispatch a notification
	if out.Following {
		// TODO: notify followw
	}

	return out, nil
}
