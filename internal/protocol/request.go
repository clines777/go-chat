package protocol

import _ "github.com/go-playground/validator/v10"

type GetTokenReq struct {
	Username string `json:"username"`
}

type GetUserInfoReq struct {
	UserId int `form:"user_id"`
}

type GetGroupInfoReq struct {
	GroupID int `form:"group_id"`
}

type ApiTokenInfo struct {
	UserID int `json:"user_id"`
}

type JoinGroupReq struct {
	GroupID int `json:"group_id" binding:"required"`
}

type LeaveGroupReq struct {
	GroupID int `json:"group_id" binding:"required"`
}

type CreateGroupReq struct {
	Title     string `json:"title" binding:"required"`
	Code      string `json:"code"`
	UserLimit int    `json:"user_limit"`
	Bulletin  string `json:"bulletin"`
	OpenJoin  bool   `json:"open_join"`
}

type SetAvatarReq struct {
	AvatarId int `json:"avatar_id" binding:"required"`
}

type ResumeReq struct {
	Token string `json:"token"`
}

type ResumeTokenInfo struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
}

type UpdateLastReadReq struct {
	GroupID int `json:"group_id"`
	ChatID  int `json:"chat_id"`
}

type UpdateGroupReq struct {
	GroupID  int    `json:"group_id" validate:"required"`
	Title    string `json:"title" validate:"required"`
	Bulletin string `json:"bulletin,omitempty"`
	Remark   string `json:"remark,omitempty"`
}
