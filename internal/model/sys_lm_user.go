package model

import "time"

type SysLmUser struct {
	ID          int64     `db:"id" json:"id"`
	LmID        int       `db:"lm_id" json:"lm_id"`
	UserID      int64     `db:"user_id" json:"user_id"`
	SiteBid     string    `db:"site_bid" json:"site_bid"`
	ExtMemberID int64     `db:"ext_member_id" json:"ext_member_id"`
	ExtUsername string    `db:"ext_username" json:"ext_username"`
	IsTaken     bool      `db:"is_taken" json:"is_taken"`
	CreateTime  time.Time `db:"create_time" json:"create_time"`
}
