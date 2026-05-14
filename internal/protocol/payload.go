package protocol

import (
	_ "github.com/go-playground/validator/v10"
)

type Payload struct {
	MsgType Type   `json:"msg_type"`
	Remark  string `json:"remark,omitempty"`
	Data    []byte `json:"data,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
	Origin  *Payload
}

type Meta struct {
	RequestID string `json:"request_id,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

type LoginReq struct {
	Token     string `json:"token" validate:"required"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

type LoginResp struct {
	UserID     int64              `json:"user_id"`
	Username   string             `json:"username"`
	ApiToken   string             `json:"api_token"`
	UserGroups []DisplayUserGroup `json:"user_groups"`
}

type DisplayUserGroup struct {
	Id            int64  `db:"id"             json:"id"`
	Title         string `db:"title"          json:"title"`
	Code          string `db:"code"           json:"code"`
	OpenJoin      bool   `db:"open_join"      json:"open_join"`
	JoinUserLevel int32  `db:"join_user_level" json:"join_user_level"`
	LastMsg       string `db:"last_msg"      json:"last_msg"`
	LastMsgTime   int64  `db:"last_msg_time" json:"last_msg_time"`
}

type EnterGroupReq struct {
	GroupId int32 `json:"group_id" validate:"required"`
}

type EnterGroupResp struct {
	Title          string `json:"title" json:"title"`
	GroupId        int32  `json:"group_id" json:"group_id"`
	GroupUserCount int32  `json:"group_user_count" json:"group_user_count"`
}

type UserSelfResp struct {
	Id         int64  `json:"id"`
	Username   string `json:"username"`
	CreateTime int64  `json:"create_time"`
	Code       string `json:"code"`
	UserLevel  int32  `json:"user_level"`
}

type GroupInfoResp struct {
	Title         string `json:"title"`
	UserTotal     int32  `json:"user_total"`
	Bulletin      string `json:"bulletin"`
	OwnerUsername string `json:"owner_username"`
	Code          string `json:"code"`
	Remark        string `json:"remark"`
}

func NewErrPayload(msgType Type, remark string, origin *Payload) *Payload {
	return &Payload{MsgType: msgType, Remark: remark, Origin: origin}
}
