package model

type ChatRecord struct {
	ID         int    `db:"id" json:"id"`
	GroupID    int    `db:"group_id" json:"group_id"`
	UserID     int    `db:"user_id" json:"user_id"`
	Type       int16  `db:"type" json:"type"`
	CreateTime int    `db:"create_time" json:"create_time"`
	UpdateTime int    `db:"update_time" json:"update_time"`
	Content    string `db:"content" json:"content"`
	Deleted    bool   `db:"deleted" json:"deleted"`
}
