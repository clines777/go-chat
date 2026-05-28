package model

type GroupUser struct {
	ID                int  `db:"id" json:"id"`
	UserID            int  `db:"user_id" json:"user_id"`
	GroupID           int  `db:"group_id" json:"group_id"`
	IsBan             bool `db:"is_ban" json:"is_ban"`
	JoinTime          int  `db:"join_time" json:"join_time"`
	RoleType          int8 `db:"role_type" json:"role_type"`
	LastVisibleChatID int  `db:"last_visible_chat_id" json:"last_visible_chat_id"`
	LastReadChatID    int  `db:"last_read_chat_id" json:"last_read_chat_id"`
	UpdateTime        int  `db:"update_time" json:"update_time"`
	Deleted           bool `db:"deleted" json:"deleted"`
}
