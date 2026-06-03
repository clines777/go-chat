package protocol

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	_ "github.com/go-playground/validator/v10"
)

// ApiResponse http 返回 response
type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func NewApiResponse(code int, msg string, data any) *ApiResponse {
	return &ApiResponse{Code: code, Message: msg, Data: data}
}

func (r *ApiResponse) H() map[string]any {
	return gin.H{"code": r.Code, "msg": r.Message, "data": r.Data}
}

//ws payload

type Payload struct {
	MsgType MsgType         `json:"msg_type"`
	ErrCode int             `json:"err_code,omitempty"`
	Remark  string          `json:"remark,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
	Meta    *Meta           `json:"meta,omitempty"`
	Origin  *Payload        `json:"-"`
}

type Meta struct {
	RequestID string `json:"request_id,omitempty"`
	Timestamp int    `json:"timestamp,omitempty"`
}

type DisplayUserGroup struct {
	Id          int    `db:"id"           json:"id"`
	Title       string `db:"title"        json:"title"`
	Code        string `db:"code"         json:"code"`
	OpenJoin    bool   `db:"open_join"    json:"open_join"`
	LastMsg     string `db:"last_msg"     json:"last_msg"`
	LastMsgTime int    `db:"last_msg_time" json:"last_msg_time"`
	CoverURL    string `db:"cover_url"    json:"cover_url"`
}

type EnterLobbyReq struct {
	Page int `json:"page"`
}

type DisplayLobbyGroup struct {
	Id          int    `db:"id"           json:"id"`
	Title       string `db:"title"        json:"title"`
	MemberCount int    `db:"member_count" json:"member_count"`
	CoverURL    string `db:"cover_url"    json:"cover_url"`
}

type EnterLobbyResp struct {
	Page   int                 `json:"page"`
	Groups []DisplayLobbyGroup `json:"groups"`
}

type EnterMyGroupReq struct {
	Page int `json:"page"`
}

type EnterMyGroupResp struct {
	Page   int                `json:"page"`
	Groups []DisplayUserGroup `json:"groups"`
}

type SendChatReq struct {
	GroupId int    `json:"group_id"`
	Content string `json:"content"`
}

type SendChatResp struct {
	Id         int `json:"id"`
	GroupID    int `json:"group_id"`
	CreateTime int `json:"create_time"`
}

type CastChatEvent struct {
	Id         int    `json:"id"`
	GroupID    int    `json:"group_id"`
	UserId     int    `json:"user_id"`
	Username   string `json:"username"`
	AvatarURL  string `json:"avatar_url"`
	Content    string `json:"content"`
	CreateTime int    `json:"create_time"`
}

type LeaveGroupEvent struct {
	GroupId int   `json:"group_id"`
	UserId  int   `json:"user_id"`
	Time    int64 `json:"time"`
}

type EnterGroupReq struct {
	GroupID int `json:"group_id" validate:"required"`
}

type ChatInfo struct {
	ID         int    `db:"id"`
	UserID     int    `db:"user_id"`
	Username   string `db:"username"`
	AvatarURL  string `db:"avatar_url" json:"avatar_url"`
	Content    string `db:"content"`
	CreateTime int    `db:"create_time"`
}

type EnterGroupResp struct {
	Title          string     `json:"title"`
	GroupId        int        `json:"group_id"`
	GroupUserCount int        `json:"group_user_count"`
	Chats          []ChatInfo `json:"chats"`
}

type UserSelfResp struct {
	Id         int    `json:"id"`
	Username   string `json:"username"`
	Nickname   string `json:"nickname"`
	AvatarURL  string `json:"avatar_url"`
	CreateTime int    `json:"create_time"`
	Code       string `json:"code"`
}

type GroupInfoResp struct {
	Title         string `json:"title"`
	UserTotal     int    `json:"user_total"`
	Bulletin      string `json:"bulletin"`
	OwnerUsername string `json:"owner_username"`
	OwnerUserID   int    `json:"owner_user_id"`
	Code          string `json:"code"`
	Remark        string `json:"remark"`
	CoverURL      string `json:"cover_url"`
}

type UpdateGroupCastEvent struct {
	GroupID  int    `json:"group_id"`
	Title    string `json:"title"`
	Bulletin string `json:"bulletin"`
	Remark   string `json:"remark"`
}

type PinChatReq struct {
	GroupID int `json:"group_id" validate:"required"`
	ChatID  int `json:"chat_id" validate:"required"`
}

type UnpinChatReq struct {
	GroupID int `json:"group_id" validate:"required"`
}

type DelChatReq struct {
	GroupID int `json:"group_id" validate:"required"`
	ChatID  int `json:"chat_id" validate:"required"`
}

type PinChatCastEvent struct {
	GroupID int `json:"group_id"`
	ChatID  int `json:"chat_id"`
}

type UnpinChatCastEvent struct {
	GroupID int `json:"group_id"`
}

type DelChatCastEvent struct {
	GroupID int `json:"group_id"`
	ChatID  int `json:"chat_id"`
}

func NewErrPayload(errCode int, remark string, origin *Payload) *Payload {
	return &Payload{MsgType: Error, ErrCode: errCode, Remark: remark, Origin: origin}
}

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

type LoginReq struct {
	Token     string `json:"token" validate:"required"`
	Timestamp int    `json:"timestamp,omitempty"`
}

type LoginResp struct {
	UserID      int                `json:"user_id"`
	Username    string             `json:"username"`
	ApiToken    string             `json:"api_token"`
	ResumeToken string             `json:"resume_token"`
	UserGroups  []DisplayUserGroup `json:"user_groups"`
}
