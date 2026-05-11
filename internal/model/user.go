package model

import "time"

type User struct {
	ID            int64     `db:"id" json:"id"`
	SiteBid       string    `db:"site_bid" json:"site_bid"`
	ExtMemberID   int64     `db:"ext_member_id" json:"ext_member_id"`
	ExtUsername   string    `db:"ext_username" json:"ext_username"`
	LastLoginTime time.Time `db:"last_login_time" json:"last_login_time"`
	AvatarID      int32     `db:"avatar_id" json:"avatar_id"`
	IsSuspended   bool      `db:"is_suspended" json:"is_suspended"`
	Code          string    `db:"code" json:"code"`
	CreateTime    time.Time `db:"create_time" json:"create_time"`
	UpdateTime    time.Time `db:"update_time" json:"update_time"`
}
