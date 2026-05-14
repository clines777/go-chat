package model

type ChatRecord struct {
	ID         int64  `db:"id" json:"id"`
	SiteBid    string `db:"site_bid" json:"site_bid"`
	GroupID    int    `db:"group_id" json:"group_id"`
	UserID     int    `db:"user_id" json:"user_id"`
	Type       int16  `db:"type" json:"type"`
	CreateTime int64  `db:"create_time" json:"create_time"`
	UpdateTime int64  `db:"update_time" json:"update_time"`
	Content    string `db:"content" json:"content"`

	Extra   []byte `db:"extra" json:"extra"`
	Deleted bool   `db:"deleted" json:"deleted"`
}
