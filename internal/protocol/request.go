package protocol

import _ "github.com/go-playground/validator/v10"

type GetTokenReq struct {
	SiteBid   string `json:"site_bid"`
	Username  string `json:"username"`
	MemberId  int64  `json:"member_id"`
	UserLevel int32  `json:"user_level,omitempty"`
}

type GetUserInfoReq struct {
	UserId int64 `form:"user_id"`
}

type ApiToken struct {
	ApiToken string `json:"api_token"`
}

type JoinGroupReq struct {
	UserID  int64 `json:"user_id"  binding:"required"`
	GroupID int64 `json:"group_id" binding:"required"`
}

type CreateGroupReq struct {
	UserID    int64  `json:"user_id"   binding:"required"`
	Title     string `json:"title"     binding:"required"`
	Code      string `json:"code"`
	UserLimit int    `json:"user_limit"`
	Bulletin  string `json:"bulletin"`
	OpenJoin  bool   `json:"open_join"`
}
