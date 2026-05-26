package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"

	"gochat/internal/protocol"
	"gochat/internal/session"
	"gochat/internal/user"
	"gochat/internal/ws"
)

func EnterSelf(ctx *ws.Ctx) *protocol.Payload {
	sess := session.Get(ctx.Client.ConnID)

	u, err := user.FindByID(sess.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return protocol.NewErrPayload(protocol.ErrUserNotFound, "user not found", ctx.Payload)
		}
		log.Printf("[EnterSelf] FindByID user=%d error: %v", sess.UserID, err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	avatarURL := ""
	if u.AvatarFilename != "" {
		avatarURL = "/static/avatars/" + u.AvatarFilename
	}

	respData, err := json.Marshal(&protocol.UserSelfResp{
		Id:         u.ID,
		Username:   u.ExtUsername,
		Nickname:   u.Nickname,
		AvatarURL:  avatarURL,
		CreateTime: u.CreateTime,
		Code:       u.Code,
	})
	if err != nil {
		log.Printf("[EnterSelf] marshal error: %v", err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	return &protocol.Payload{MsgType: protocol.EnterSelfOk, Data: respData, Meta: ctx.Payload.Meta}
}
