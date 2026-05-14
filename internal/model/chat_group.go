package model

type ChatGroup struct {
	ID                    int64  `db:"id" json:"id"`
	SiteBid               string `db:"site_bid" json:"site_bid"`
	IsDismiss             bool   `db:"is_dismiss" json:"is_dismiss"`
	Sort                  int    `db:"sort" json:"sort"`
	Visible               bool   `db:"visible" json:"visible"`
	Code                  string `db:"code" json:"code"`
	OpenJoin              bool   `db:"open_join" json:"open_join"`
	UserLimit             int    `db:"user_limit" json:"user_limit"`
	Bulletin              string `db:"bulletin" json:"bulletin"`
	AllowURL              bool   `db:"allow_url" json:"allow_url"`
	PinChatID             int64  `db:"pin_chat_id" json:"pin_chat_id"`
	SpeakUserLevel        int    `db:"speak_user_level" json:"speak_user_level"`
	JoinUserLevel         int    `db:"join_user_level" json:"join_user_level"`
	OwnerExtMemberID      int    `db:"owner_ext_member_id" json:"owner_ext_member_id"`
	OwnerExtUsername      string `db:"owner_ext_username" json:"owner_ext_username"`
	OwnerUserID           int    `db:"owner_user_id" json:"owner_user_id"`
	ModExtMemberIDs       string `db:"mod_ext_member_ids" json:"mod_ext_member_ids"`
	ModExtUsernames       string `db:"mod_ext_usernames" json:"mod_ext_usernames"`
	ModUserIDs            string `db:"mod_user_ids" json:"mod_user_ids"`
	Title                 string `db:"title" json:"title"`
	Remark                string `db:"remark" json:"remark"`
	CreateTime            int64  `db:"create_time" json:"create_time"`
	UpdateTime            int64  `db:"update_time" json:"update_time"`
	LobbyCoverPicURL      string `db:"lobby_cover_pic_url" json:"lobby_cover_pic_url"`
	MyGroupCoverPicURL    string `db:"my_group_cover_pic_url" json:"my_group_cover_pic_url"`
	ExtLobbyCoverPicURL   string `db:"ext_lobby_cover_pic_url" json:"ext_lobby_cover_pic_url"`
	ExtMyGroupCoverPicURL string `db:"ext_my_group_cover_pic_url" json:"ext_my_group_cover_pic_url"`
}
