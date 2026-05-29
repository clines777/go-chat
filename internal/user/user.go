package user

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"gochat/internal/infra/db"
	"gochat/internal/infra/redis"
	"gochat/internal/model"
	"gochat/internal/protocol"
	"time"
)

const ResumeTokenTTL = 60 * 60 * 1

func GenerateResumeToken(u *model.User) (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)
	tokenInfo := protocol.ResumeTokenInfo{
		UserID:   u.ID,
		Username: u.Username,
	}
	if err := redis.GetRedis().SetJSON(protocol.ResumeTokenKey(token), tokenInfo, ResumeTokenTTL); err != nil {
		return "", err
	}
	return token, nil
}

func GetResumeToken(token string) (*protocol.ResumeTokenInfo, error) {
	var payload protocol.ResumeTokenInfo
	ok, err := redis.GetRedis().GetJSON(protocol.ResumeTokenKey(token), &payload)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return &payload, nil
}

func RefreshResumeToken(token string) {
	_ = redis.GetRedis().Expire(protocol.ResumeTokenKey(token), ResumeTokenTTL)
}

func GenerateApiToken(userID int) (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)
	payload := protocol.ApiTokenPayload{UserID: userID}
	if err := redis.GetRedis().SetJSON(protocol.ApiTokenKey(token), payload, 24*time.Hour); err != nil {
		return "", err
	}
	return token, nil
}

func GetLoginToken(req protocol.LoginReq) (*protocol.GetTokenReq, error) {
	r := redis.GetRedis()
	key := protocol.LoginTokenKey(req.Token)

	var tokenUser protocol.GetTokenReq
	ok, err := r.GetJSON(key, &tokenUser)
	if err != nil || !ok {
		return nil, err
	}

	_ = r.Del(key)

	return &tokenUser, nil
}

func Login(tokenInfo *protocol.GetTokenReq) (*model.User, error) {
	user, err := findUser(tokenInfo.Username)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		userCode := genUserCode(tokenInfo.Username)
		user, err = createUser(tokenInfo, userCode)
		if err != nil {
			return nil, err
		}
	}

	user, err = updateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func findUser(username string) (*model.User, error) {
	d, err := db.GetDBConn()
	if err != nil {
		return nil, err
	}

	findSql, args, _ := d.Builder.
		Select(
			"u.id", "u.username", "u.code",
			"u.is_suspended", "u.avatar_id",
			"COALESCE(av.filename, '') AS avatar_filename",
		).
		From(`"user" u`).
		LeftJoin("avatar av ON av.id = u.avatar_id").
		Where(sq.Eq{"u.username": username}).
		Limit(1).ToSql()

	var u model.User
	if err := d.DB.Get(&u, findSql, args...); err != nil {
		return nil, err
	}

	return &u, nil
}

func FindByID(userID int) (*model.User, error) {
	d, err := db.GetDBConn()
	if err != nil {
		return nil, err
	}

	findSql, args, _ := d.Builder.
		Select(
			"u.id", "u.username", "u.nickname",
			"u.code", "u.is_suspended", "u.avatar_id", "u.create_time",
			"COALESCE(av.filename, '') AS avatar_filename",
		).
		From(`"user" u`).
		LeftJoin("avatar av ON av.id = u.avatar_id").
		Where(sq.Eq{"u.id": userID}).
		Limit(1).ToSql()

	var m model.User
	if err := d.DB.Get(&m, findSql, args...); err != nil {
		return nil, err
	}

	return &m, nil
}

func createUser(tokenInfo *protocol.GetTokenReq, userCode string) (*model.User, error) {
	d, err := db.GetDBConn()
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	var id int
	err = d.Builder.
		Insert(`"user"`).
		Columns("username", "code", "last_login_time", "create_time", "update_time").
		Values(tokenInfo.Username, userCode, now, now, now).
		Suffix("RETURNING id").
		RunWith(d.DB).
		QueryRow().
		Scan(&id)
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:         id,
		Username:   tokenInfo.Username,
		Code:       userCode,
		CreateTime: int(now),
		UpdateTime: int(now),
	}, nil
}

func UpdateAvatar(userID int, avatarID int) error {
	d, err := db.GetDBConn()
	if err != nil {
		return err
	}

	_, err = d.Builder.
		Update(`"user"`).
		Set("avatar_id", avatarID).
		Where(sq.Eq{"id": userID}).
		RunWith(d.DB).Exec()
	return err
}

func updateUser(u *model.User) (*model.User, error) {
	d, err := db.GetDBConn()
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	_, err = d.Builder.
		Update(`"user"`).
		Set("last_login_time", now).
		Where(sq.Eq{"id": u.ID}).
		RunWith(d.DB).
		Exec()

	if err != nil {
		return nil, err
	}

	u.LastLoginTime = int(now)

	return u, nil
}
