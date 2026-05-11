package user

import (
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"gochat/internal/infra/db"
	"gochat/internal/infra/redis"
	"gochat/internal/model"
	"gochat/internal/protocol"
	"gochat/internal/ws"
	"strings"
	"time"
)

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

	sess := &ws.Session{
		ConnID:  c.Client.ConnID,
		UserID:  user.ID,
		SiteBid: tokenInfo.SiteBid,
	}

	r := redis.GetRedis()
	if err := r.SetJSON(protocol.SessionKey(c.Client.ConnID), sess, 24*time.Hour); err != nil {
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
		Select("*").
		From("user").
		Where(sq.Eq{"site_bid": siteBid, "ext_username": username}).
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
		Insert("member").
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

func updateUser(u *model.User) (*model.User, error) {
	d, err := db.GetDBConn()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	_, err = d.Builder.
		Update("member").
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
