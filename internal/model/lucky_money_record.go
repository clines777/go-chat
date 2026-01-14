package model

import "time"

type LuckyMoneyRecord struct {
	ID              int64  `db:"id" json:"id"`
	LmID            int    `db:"lm_id" json:"lm_id"`
	UserID          int64  `db:"user_id" json:"user_id"`
	GroupID         int    `db:"group_id" json:"group_id"`
	SiteBid         string `db:"site_bid" json:"site_bid"`
	ExtMemberID     int64  `db:"ext_member_id" json:"ext_member_id"`
	ExtUsername     string `db:"ext_username" json:"ext_username"`
	ExtPlatformType int16  `db:"ext_platform_type" json:"ext_platform_type"`
	PackID          int    `db:"pack_id" json:"pack_id"`
	SendState       int    `db:"send_state" json:"send_state"`
	Device          int16  `db:"device" json:"device"`
	Tpl             int16  `db:"tpl" json:"tpl"`

	// 金額相關（numeric）
	Amount        float64 `db:"amount" json:"amount"`
	RequireAmount float64 `db:"require_amount" json:"require_amount"`

	CreateTime time.Time `db:"create_time" json:"create_time"`
	FundTime   time.Time `db:"fund_time" json:"fund_time"`
}
