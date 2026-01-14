package model

import "time"

type Member struct {
	ID                int64     `db:"id" json:"id"`
	SiteBid           string    `db:"site_bid" json:"site_bid"`
	ExtMemberID       int       `db:"ext_member_id" json:"ext_member_id"`
	ExtUsername       string    `db:"ext_username" json:"ext_username"`
	ExtPlatformType   int16     `db:"ext_platform_type" json:"ext_platform_type"`
	UserLevel         int       `db:"user_level" json:"user_level"`
	LastLoginTime     time.Time `db:"last_login_time" json:"last_login_time"`
	AvatarID          int       `db:"avatar_id" json:"avatar_id"`
	IsGlobalBan       bool      `db:"is_global_ban" json:"is_global_ban"`
	Code              string    `db:"code" json:"code"`
	SysChatLastReadID int       `db:"sys_chat_last_read_id" json:"sys_chat_last_read_id"`
	CreateTime        time.Time `db:"create_time" json:"create_time"`
	UpdateTime        time.Time `db:"update_time" json:"update_time"`
}
