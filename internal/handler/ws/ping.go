package handler

import (
	"gochat/internal/ws"
)

func Ping(ctx *ws.Ctx) ([]byte, error) {
	return []byte("pong"), nil
}
