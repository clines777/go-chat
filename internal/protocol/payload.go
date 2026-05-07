package protocol

import _ "github.com/go-playground/validator/v10"

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

func NewErrPayload(msgType Type, remark string, origin *Payload) *Payload {
	return &Payload{MsgType: msgType, Remark: remark, Origin: origin}
}
