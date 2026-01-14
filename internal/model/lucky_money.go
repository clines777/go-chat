package model

import "time"

type LuckyMoney struct {
	ID          int64  `db:"id" json:"id"`
	GroupID     int    `db:"group_id" json:"group_id"`
	SiteBid     string `db:"site_bid" json:"site_bid"`
	UserID      int64  `db:"user_id" json:"user_id"`
	ExtMemberID int64  `db:"ext_member_id" json:"ext_member_id"`
	ExtUsername string `db:"ext_username" json:"ext_username"`
	AdminID     string `db:"admin_id" json:"admin_id"`

	// 金額相關（numeric）
	UnitAmount   float64 `db:"unit_amount" json:"unit_amount"`
	Count        int     `db:"count" json:"count"`
	TotalAmount  float64 `db:"total_amount" json:"total_amount"`
	RequireMulti float64 `db:"require_multi" json:"require_multi"`

	// 類型 / 狀態
	Title         string `db:"title" json:"title"`
	GameCompanyID int    `db:"game_company_id" json:"game_company_id"`
	SceneType     int16  `db:"scene_type" json:"scene_type"`
	CreateType    int16  `db:"create_type" json:"create_type"`
	AmountType    int16  `db:"amount_type" json:"amount_type"`
	BonusType     int16  `db:"bonus_type" json:"bonus_type"`
	CloseType     int16  `db:"close_type" json:"close_type"`

	// 時間
	StartTime  time.Time `db:"start_time" json:"start_time"`
	EndTime    time.Time `db:"end_time" json:"end_time"`
	CreateTime time.Time `db:"create_time" json:"create_time"`
	UpdateTime time.Time `db:"update_time" json:"update_time"`
}
