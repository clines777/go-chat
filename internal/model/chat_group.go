package model

type ChatGroup struct {
	ID             int64  `db:"id" json:"id"`
	IsDismiss      bool   `db:"is_dismiss" json:"is_dismiss"`
	Sort           int    `db:"sort" json:"sort"`
	Visible        bool   `db:"visible" json:"visible"`
	Code           string `db:"code" json:"code"`
	OpenJoin       bool   `db:"open_join" json:"open_join"`
	UserLimit      int    `db:"user_limit" json:"user_limit"`
	Bulletin       string `db:"bulletin" json:"bulletin"`
	AllowURL       bool   `db:"allow_url" json:"allow_url"`
	PinChatID      int64  `db:"pin_chat_id" json:"pin_chat_id"`
	OwnerUserName  string `db:"owner_username" json:"owner_username"`
	OwnerUserID    int    `db:"owner_user_id" json:"owner_user_id"`
	Title          string `db:"title" json:"title"`
	Remark         string `db:"remark" json:"remark"`
	CreateTime     int64  `db:"create_time" json:"create_time"`
	UpdateTime     int64  `db:"update_time" json:"update_time"`
	CoverFileName  string `db:"cover_filename" json:"cover_filename"`
}
