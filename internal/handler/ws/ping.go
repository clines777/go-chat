package handler

import (
	"gochat/internal/protocol"
	"gochat/internal/ws"
)

func Ping(ctx *ws.Ctx) *protocol.Payload {
	return &protocol.Payload{MsgType: protocol.Pong}
}
