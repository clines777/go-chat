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
	"gochat/internal/session"
	"gochat/internal/ws"
	"strings"
	"time"
)

const resumeTokenTTL = 7 * 24 * time.Hour

func GenerateResumeToken(u *model.User) (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)
	payload := protocol.ResumeTokenPayload{
		UserID:   u.ID,
		SiteBid:  u.SiteBid,
		Username: u.ExtUsername,
	}
	if err := redis.GetRedis().SetJSON(protocol.ResumeTokenKey(token), payload, resumeTokenTTL); err != nil {
		return "", err
	}
	return token, nil
}

func GetResumeToken(token string) (*protocol.ResumeTokenPayload, error) {
	var payload protocol.ResumeTokenPayload
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
	_ = redis.GetRedis().Expire(protocol.ResumeTokenKey(token), resumeTokenTTL)
}

func GenerateApiToken(userID int64, siteBid string) (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)
	payload := protocol.ApiTokenPayload{UserID: userID, SiteBid: siteBid}
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

	tokenUser.SiteBid = strings.ToUpper(tokenUser.SiteBid)

	_ = r.Del(key)

	return &tokenUser, nil
}

func Login(c *ws.Ctx, tokenInfo *protocol.GetTokenReq) (*model.User, error) {

	user, err := findUser(tokenInfo.SiteBid, tokenInfo.Username)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		userCode := genUserCode(tokenInfo.SiteBid, tokenInfo.MemberId, tokenInfo.Username)
		user, err = createUser(tokenInfo, userCode)
		if err != nil {
			return nil, err
		}
	}

	user, err = updateUser(user)
	if err != nil {
		return nil, err
	}

	sess := &session.Session{
		ConnID:    c.Client.ConnID,
		UserID:    user.ID,
		SiteBid:   user.SiteBid,
		Username:  user.ExtUsername,
		UserLevel: user.UserLevel,
		InGroupID: 0,
		Scene:     protocol.SceneMyGroup,
	}

	err = session.Set(c.Client.ConnID, sess)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func findUser(siteBid string, username string) (*model.User, error) {
	d, err := db.GetDBConn()
	if err != nil {
		return nil, err
	}

	findSql, args, _ := d.Builder.
		Select(
			"u.id", "u.site_bid", "u.ext_member_id", "u.ext_username", "u.code",
			"u.user_level", "u.is_suspended", "u.avatar_id",
			"COALESCE(av.filename, '') AS avatar_filename",
		).
		From(`"user" u`).
		LeftJoin("avatar av ON av.id = u.avatar_id").
		Where(sq.Eq{"u.site_bid": siteBid, "u.ext_username": username}).
		Limit(1).ToSql()

	var u model.User
	if err := d.DB.Get(&u, findSql, args...); err != nil {
		return nil, err
	}

	return &u, nil
}

var userColumns = []string{"id", "site_bid", "ext_member_id", "ext_username", "nickname", "code", "user_level", "is_suspended", "avatar_id", "create_time"}

func FindByID(userID int64) (*model.User, error) {
	d, err := db.GetDBConn()
	if err != nil {
		return nil, err
	}

	findSql, args, _ := d.Builder.
		Select(
			"u.id", "u.site_bid", "u.ext_member_id", "u.ext_username", "u.nickname",
			"u.code", "u.user_level", "u.is_suspended", "u.avatar_id", "u.create_time",
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

func FindByExtMember(siteBid string, extMemberID int64) (*model.User, error) {
	d, err := db.GetDBConn()
	if err != nil {
		return nil, err
	}

	findSql, args, _ := d.Builder.
		Select(userColumns...).
		From(`"user"`).
		Where(sq.Eq{"site_bid": siteBid, "ext_member_id": extMemberID}).
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

	now := time.Now()
	insertSql, insertArgs, _ := d.Builder.
		Insert(`"user"`).
		Columns("site_bid", "ext_member_id", "ext_username", "code", "last_login_time", "create_time", "update_time").
		Values(tokenInfo.SiteBid, tokenInfo.MemberId, tokenInfo.Username, userCode, now, now, now).
		Suffix("RETURNING *").
		ToSql()

	var u model.User
	if err := d.DB.Get(&u, insertSql, insertArgs...); err != nil {
		return nil, err
	}
	return &u, nil
}

func UpdateAvatar(userID int64, avatarID int32) error {
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

	u.LastLoginTime = now

	return u, nil
}
