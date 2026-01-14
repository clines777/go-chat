package model

import "time"

type LuckyMoneyPack struct {
	ID   int64 `db:"id" json:"id"`
	LmID int   `db:"lm_id" json:"lm_id"`
	No   int   `db:"no" json:"no"`

	// 金額（numeric 13,3）
	Amount float64 `db:"amount" json:"amount"`

	UserID     int64     `db:"user_id" json:"user_id"`
	IsTaken    bool      `db:"is_taken" json:"is_taken"`
	CreateTime time.Time `db:"create_time" json:"create_time"`
	UpdateTime time.Time `db:"update_time" json:"update_time"`
}
