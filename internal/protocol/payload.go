package protocol

import (
	"encoding/json"
	_ "github.com/go-playground/validator/v10"
)

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

func NewErrPayload(errCode int, remark string, origin *Payload) *Payload {
	return &Payload{MsgType: Error, ErrCode: errCode, Remark: remark, Origin: origin}
}
