package model

import "time"

type GroupQuotaLog struct {
	ID           int64  `db:"id" json:"id"`
	SiteBid      string `db:"site_bid" json:"site_bid"`
	GroupID      int    `db:"group_id" json:"group_id"`
	AlterType    int16  `db:"alter_type" json:"alter_type"`
	QuotaType    int16  `db:"quota_type" json:"quota_type"`
	ExtUsername  string `db:"ext_username" json:"ext_username"`
	AdminID      int64  `db:"admin_id" json:"admin_id"`
	ActionType   int16  `db:"action_type" json:"action_type"`
	OperatorRole int16  `db:"operator_role" json:"operator_role"`

	// 金額欄位（numeric 13,3）
	Amount     float64   `db:"amount" json:"amount"`
	CreateTime time.Time `db:"create_time" json:"create_time"`
}
