package model

import "time"

type ChatRecord struct {
	ID         int64     `db:"id" json:"id"`
	SiteBid    string    `db:"site_bid" json:"site_bid"`
	GroupID    int       `db:"group_id" json:"group_id"`
	UserID     int       `db:"user_id" json:"user_id"`
	Type       int16     `db:"type" json:"type"`
	CreateTime time.Time `db:"create_time" json:"create_time"`
	UpdateTime time.Time `db:"update_time" json:"update_time"`
	Content    string    `db:"content" json:"content"`
	BadWords   string    `db:"bad_words" json:"bad_words"`
	CustomID   int       `db:"custom_id" json:"custom_id"`

	// jsonb：DB 層用 Raw bytes 最穩定
	Extra   []byte `db:"extra" json:"extra"`
	Deleted bool   `db:"deleted" json:"deleted"`
}
