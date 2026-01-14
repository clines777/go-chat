package model

import "time"

type SysChatAudience struct {
	ID          int64     `db:"id" json:"id"`
	UserID      int64     `db:"user_id" json:"user_id"`
	ExtMemberID int64     `db:"ext_member_id" json:"ext_member_id"`
	ExtUsername string    `db:"ext_username" json:"ext_username"`
	ChatID      int       `db:"chat_id" json:"chat_id"`
	CreateTime  time.Time `db:"create_time" json:"create_time"`
}
