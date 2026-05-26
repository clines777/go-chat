package protocol

import _ "github.com/go-playground/validator/v10"

type GetTokenReq struct {
	Username string `json:"username"`
}

type GetUserInfoReq struct {
	UserId int64 `form:"user_id"`
}

type GetGroupInfoReq struct {
	GroupID int64 `form:"group_id"`
}

type ApiTokenPayload struct {
	UserID int64 `json:"user_id"`
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

type ResumeReq struct {
	Token string `json:"token"`
}

type ResumeTokenPayload struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
}
