package model

type User struct {
	ID             int    `db:"id" json:"id"`
	Username       string `db:"username" json:"username"`
	Nickname       string `db:"nickname" json:"nickname"`
	LastLoginTime  int    `db:"last_login_time" json:"last_login_time"`
	AvatarID       int    `db:"avatar_id" json:"avatar_id"`
	IsSuspended    bool   `db:"is_suspended" json:"is_suspended"`
	Code           string `db:"code" json:"code"`
	CreateTime     int    `db:"create_time" json:"create_time"`
	UpdateTime     int    `db:"update_time" json:"update_time"`
	AvatarFilename string `db:"avatar_filename" json:"-"`
}
