package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/disintegration/imaging"
	gonanoid "github.com/matoous/go-nanoid"
)

// 아바타 용량 제한
// MaxAvatarBytes to read
const MaxAvatarBytes = 5 << 20 //5MB

var (
	//이메일 및 이름 확인(정규식 사용)
	rxEmail    = regexp.MustCompile("^[^\\s@]+@[^\\s@]+\\.[^\\s@]+$")
	rxUsername = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_-]{0,17}$")
	avatarsDir = path.Join("web", "static", "img", "avatars")
)

var (
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
	// ErrUnsupportedAvatarFormat used for unsupported avatar format.
	ErrUnsupportedAvatarFormat = errors.New("only png and jpeg allowed as avatar")
)

//User Model
type User struct {
	ID        int64   `json:"id,omitempty"` //omitempty는 필드에서 값 반환 금지
	UserName  string  `json:"user_name"`
	AvatarURL *string `json:"avatarUrl"`
}

//디테일한 유저 구조체
//Userprofile model.
type UserProfile struct {
	User
	Email          string `json:"email,omitempty"`
	FollowersCount int    `json:"followers_count"`
	FolloweesCount int    `json:"followees_count"`
	Me             bool   `json:"me"`
	Following      bool   `json:"following"`
	Followeed      bool   `json:"followeed"`
}

//팔로워 카운트
//ToggleFollowOutput response
type ToggleFollowOutput struct {
	Following      bool `json:"following"`
	FollowersCount int  `json:"followers_count"`
}

// 유저 생성
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

//유저 이름을 오름차순으로 하여 유저 리스트 생성 및 유저 검색
//Users in ascending order with forward pagination and filtered by username
func (s *Service) Users(ctx context.Context, search string, first int, after string) ([]UserProfile, error) {
	search = strings.TrimSpace(search)
	first = normailizePageSize(first)
	after = strings.TrimSpace(after)
	uid, auth := ctx.Value(KeyAuthUserID).(int64)
	//인증된 유저의 요청만 처리하도록 하는 기능
	//USer()와 다른 방식 - go template
	query, args, err := buildQuery(`
	SELECT id, email, username, avatar followers_count, followees_count
	{{if .auth}}
	, followers.follower_id IS NOT NULL AS following
	, followees.followee_id IS NOT NULL AS followeed
	{{end}}
	FROM users
	{{if .auth}}
	LEFT JOIN follows AS followers ON followers.follower_id = @uid AND followers.followee_id = users.id
	LEFT JOIN follows AS followees ON followees.follower_id = users.id AND followees.followee_id = @uid
	{{end}}
	{{if or .search .after}} WHERE {{end}}
	{{if .search}} username ILIKE '%' || @search || '%'{{end}}
	{{if and .search .after}}AND{{end}}
	{{if .after}}username > @after{{end}}
	ORDER BY username ASC
	LIMIT @first`, map[string]interface{}{
		"auth":   auth,
		"uid":    uid,
		"search": search,
		"first":  first,
		"after":  after,
	})
	if err != nil {
		return nil, fmt.Errorf("could not build users sql query: %v", err)
	}

	log.Printf("users query: %s\nargs: %v\n", query, args)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("could not query select users: %v", err)
	}

	defer rows.Close()
	uu := make([]UserProfile, 0, first)
	for rows.Next() {
		var u UserProfile
		var avatar sql.NullString
		dest := []interface{}{&u.ID, &u.Email, &u.UserName, &avatar, &u.FollowersCount, &u.FolloweesCount}
		if auth {
			dest = append(dest, &u.Following, &u.Followeed)
		}
		if err = rows.Scan(dest...); err != nil {
			return nil, fmt.Errorf("could not scan user: %v", err)
		}

		u.Me = auth && uid == u.ID
		if !u.Me {
			u.ID = 0
			u.Email = ""
		}
		if avatar.Valid {
			avatarURL := s.origin + "/img/avatars/" + avatar.String
			u.AvatarURL = &avatarURL
		}
		uu = append(uu, u)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("could not iterate user rows: %v", err)
	}

	return uu, nil
}

func (s *Service) userByID(ctx context.Context, id int64) (User, error) {
	//Database query with the user ID in the context
	var u User
	var avatar sql.NullString
	query := "SELECT username, avatar FROM users WHERE id = $1"
	err := s.db.QueryRowContext(ctx, query, id).Scan(&u.UserName, &avatar)
	if err == sql.ErrNoRows {
		return u, ErrUserNotFound
	}

	if err != nil {
		return u, fmt.Errorf("could not query select user: %v", err)
	}

	u.ID = id
	if avatar.Valid {
		avatarURL := s.origin + "/img/avatars/" + avatar.String
		u.AvatarURL = &avatarURL
	}
	return u, nil
}

// 유저 프로필
// User selects on user from the database with the given username.
func (s *Service) User(ctx context.Context, username string) (UserProfile, error) {
	var u UserProfile

	//유저 이름 확인
	//validate the username as well
	username = strings.TrimSpace(username)
	if !rxUsername.MatchString(username) {
		return u, ErrInvalidUsername
	}
	//uid와 auth에 요청을 보낸 사람 ID 입력
	uid, auth := ctx.Value(KeyAuthUserID).(int64)

	//인증된 요청만 처리하도록 조건 설정
	//only authenticated request && query dynamically
	var avatar sql.NullString
	args := []interface{}{username}
	dest := []interface{}{&u.ID, &u.Email, &avatar, &u.FollowersCount, &u.FolloweesCount}
	query := "SELECT id, email, avatar, followers_count, followees_count "
	if auth {
		query += ", " +
			"followers.follower_id IS NOT NULL AS following, " +
			"followees.followee_id IS NOT NULL AS followeed "
		dest = append(dest, &u.Following, &u.Followeed)
	}
	query += "FROM users "
	if auth {
		query += "LEFT JOIN follows AS followers ON followers.follower_id = $2 AND followers.followee_id = users.id " +
			"LEFT JOIN follows AS followees ON followees.follower_id = users.id AND followees.followee_id = $2 "
		args = append(args, uid)
		dest = append(dest)
	}
	query += "WHERE username = $1"
	err := s.db.QueryRowContext(ctx, query, args...).Scan(dest...)

	if err == sql.ErrNoRows {
		return u, ErrUserNotFound
	}

	if err != nil {
		return u, fmt.Errorf("could not query select user: %v", err)
	}

	u.UserName = username
	u.Me = auth && uid == u.ID
	//만약 요청을 보낸 사람이 인증되지 않은 사용자이거나
	//본인이 아닐시에는 id와 email 감추기
	//reset user profile
	if !u.Me {
		u.ID = 0
		u.Email = ""
	}
	if avatar.Valid {
		avatarURL := s.origin + "/img/avatars/" + avatar.String
		u.AvatarURL = &avatarURL
	}
	return u, nil
}

// 아바타 생성
// UpdateAvatar of the authenticated user returning the new avatar URL
func (s *Service) UpdateAvatar(ctx context.Context, r io.Reader) (string, error) {
	// 유저 확인
	// checking authentication
	uid, ok := ctx.Value(KeyAuthUserID).(int64)
	if !ok {
		return "", ErrUnauthenticated
	}

	// 아바타 용량 제한, 형식 제한
	r = io.LimitReader(r, MaxAvatarBytes)
	img, format, err := image.Decode(r)
	if err != nil {
		return "", fmt.Errorf("could not read avatar: %v", err)
	}

	if format != "png" && format != "jpeg" {
		return "", ErrUnsupportedAvatarFormat
	}

	// 유저 이름에 맞는 아바타 자동 생성
	avatar, err := gonanoid.Nanoid()
	if err != nil {
		return "", fmt.Errorf("could not generate avatar filename: %v", err)
	}

	// 추가한 아바타 사진 이름에 형식 추가
	if format == "png" {
		avatar += ".png"
	} else {
		avatar += ".jpg"
	}

	// 아바타 사진이 저정된 경로 불러오기
	avatarPath := path.Join(avatarsDir, avatar)
	f, err := os.Create(avatarPath)
	if err != nil {
		return "", fmt.Errorf("could not create avatar: %v", err)
	}

	// 이미지 크기 변환
	defer f.Close()
	img = imaging.Fill(img, 400, 400, imaging.Center, imaging.CatmullRom)

	if format == "png" {
		err = png.Encode(f, img)
	} else {
		err = jpeg.Encode(f, img, nil)
	}

	if err != nil {
		return "", fmt.Errorf("could not write avatar for disk: %v", err)
	}

	var oldAvatar sql.NullString

	//새로운 아바타가 업데이트 됐을 때 기존의 아바타 사진을 자동으로 지움
	if err = s.db.QueryRowContext(ctx, `
		UPDATE users SET avatar = $1 WHERE id = $2
		RETURNING (SELECT avatar FROM users WHERE id = $2) AS old_avatar`, avatar, uid).Scan(&oldAvatar); err != nil {
		defer os.Remove(avatarPath)
		return "", fmt.Errorf("could not update avatar: %v", err)
	}

	if oldAvatar.Valid {
		defer os.Remove(path.Join(avatarsDir, oldAvatar.String))
	}

	return s.origin + "/img/avatars/" + avatar, nil
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

//팔로워 이름을 오름차순으로 하여 유저 리스트 생성 및 유저 검색
//Followers in ascending order with forward pagination and filtered by username
func (s *Service) Followers(ctx context.Context, username string, first int, after string) ([]UserProfile, error) {
	username = strings.TrimSpace(username)
	if !rxUsername.MatchString(username) {
		return nil, ErrInvalidUsername
	}
	first = normailizePageSize(first)
	after = strings.TrimSpace(after)
	uid, auth := ctx.Value(KeyAuthUserID).(int64)
	//인증된 유저의 요청만 처리하도록 하는 기능
	//USer()와 다른 방식 - go template
	query, args, err := buildQuery(`
	SELECT id, email, username, avatar, followers_count, followees_count
	{{if .auth}}
	, followers.follower_id IS NOT NULL AS following
	, followees.followee_id IS NOT NULL AS followeed
	{{end}}
	FROM follows
	INNER JOIN users ON follows.follower_id = users.id
	{{if .auth}}
	LEFT JOIN follows AS followers ON followers.follower_id = @uid AND followers.followee_id = users.id
	LEFT JOIN follows AS followees ON followees.follower_id = users.id AND followees.followee_id = @uid
	{{end}}
	WHERE follows.followee_id = (SELECT id FROM users WHERE username = @username)
	{{if .after}}AND username > @after{{end}}
	ORDER BY username ASC
	LIMIT @first`, map[string]interface{}{
		"auth":     auth,
		"uid":      uid,
		"username": username,
		"first":    first,
		"after":    after,
	})
	if err != nil {
		return nil, fmt.Errorf("could not build followers sql query: %v", err)
	}

	log.Printf("users query: %s\nargs: %v\n", query, args)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("could not query select followers: %v", err)
	}

	defer rows.Close()
	uu := make([]UserProfile, 0, first)
	for rows.Next() {
		var u UserProfile
		var avatar sql.NullString
		dest := []interface{}{&u.ID, &u.Email, &u.UserName, &avatar, &u.FollowersCount, &u.FolloweesCount}
		if auth {
			dest = append(dest, &u.Following, &u.Followeed)
		}
		if err = rows.Scan(dest...); err != nil {
			return nil, fmt.Errorf("could not scan follower: %v", err)
		}

		u.Me = auth && uid == u.ID
		if !u.Me {
			u.ID = 0
			u.Email = ""
		}
		if avatar.Valid {
			avatarURL := s.origin + "/img/avatars/" + avatar.String
			u.AvatarURL = &avatarURL
		}
		uu = append(uu, u)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("could not iterate follower rows: %v", err)
	}

	return uu, nil
}

//팔로잉한 사람 이름을 오름차순으로 하여 유저 리스트 생성 및 유저 검색
//Followees in ascending order with forward pagination and filtered by username
func (s *Service) Followees(ctx context.Context, username string, first int, after string) ([]UserProfile, error) {
	username = strings.TrimSpace(username)
	if !rxUsername.MatchString(username) {
		return nil, ErrInvalidUsername
	}
	first = normailizePageSize(first)
	after = strings.TrimSpace(after)
	uid, auth := ctx.Value(KeyAuthUserID).(int64)
	//인증된 유저의 요청만 처리하도록 하는 기능
	//USer()와 다른 방식 - go template
	query, args, err := buildQuery(`
	SELECT id, email, username, avatar followers_count, followees_count
	{{if .auth}}
	, followers.follower_id IS NOT NULL AS following
	, followees.followee_id IS NOT NULL AS followeed
	{{end}}
	FROM follows
	INNER JOIN users ON follows.followee_id = users.id
	{{if .auth}}
	LEFT JOIN follows AS followers ON followers.follower_id = @uid AND followers.followee_id = users.id
	LEFT JOIN follows AS followees ON followees.follower_id = users.id AND followees.followee_id = @uid
	{{end}}
	WHERE follows.follower_id = (SELECT id FROM users WHERE username = @username)
	{{if .after}}AND username > @after{{end}}
	ORDER BY username ASC
	LIMIT @first`, map[string]interface{}{
		"auth":     auth,
		"uid":      uid,
		"username": username,
		"first":    first,
		"after":    after,
	})
	if err != nil {
		return nil, fmt.Errorf("could not build followees sql query: %v", err)
	}

	log.Printf("users query: %s\nargs: %v\n", query, args)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("could not query select followees: %v", err)
	}

	defer rows.Close()
	uu := make([]UserProfile, 0, first)
	for rows.Next() {
		var u UserProfile
		var avatar sql.NullString
		dest := []interface{}{&u.ID, &u.Email, &u.UserName, &avatar, &u.FollowersCount, &u.FolloweesCount}
		if auth {
			dest = append(dest, &u.Following, &u.Followeed)
		}
		if err = rows.Scan(dest...); err != nil {
			return nil, fmt.Errorf("could not scan followee: %v", err)
		}

		u.Me = auth && uid == u.ID
		if !u.Me {
			u.ID = 0
			u.Email = ""
		}
		if avatar.Valid {
			avatarURL := s.origin + "/img/avatars/" + avatar.String
			u.AvatarURL = &avatarURL
		}
		uu = append(uu, u)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("could not iterate followee rows: %v", err)
	}

	return uu, nil
}
