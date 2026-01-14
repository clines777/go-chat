package model

import "time"

type SysChatContent struct {
	ID           int64     `db:"id" json:"id"`
	AudienceType int16     `db:"audience_type" json:"audience_type"`
	Deleted      bool      `db:"deleted" json:"deleted"`
	SiteBid      string    `db:"site_bid" json:"site_bid"`
	AdminID      int       `db:"admin_id" json:"admin_id"`
	ContentType  int16     `db:"content_type" json:"content_type"`
	CreateTime   time.Time `db:"create_time" json:"create_time"`
	UpdateTime   time.Time `db:"update_time" json:"update_time"`
	Content      string    `db:"content" json:"content"`
	CustomID     int       `db:"custom_id" json:"custom_id"`

	// jsonb：DB 層用 Raw bytes 最穩定
	Extra []byte `db:"extra" json:"extra"`
}
