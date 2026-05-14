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

type GetGroupInfoReq struct {
	GroupID int64 `form:"group_id"`
}

type ApiTokenPayload struct {
	UserID  int64  `json:"user_id"`
	SiteBid string `json:"site_bid"`
}

type JoinGroupReq struct {
	GroupID int64 `json:"group_id" binding:"required"`
}

type CreateGroupReq struct {
	Title     string `json:"title" binding:"required"`
	Code      string `json:"code"`
	UserLimit int    `json:"user_limit"`
	Bulletin  string `json:"bulletin"`
	OpenJoin  bool   `json:"open_join"`
}

type SetAvatarReq struct {
	AvatarId int32 `json:"avatar_id" binding:"required"`
}
