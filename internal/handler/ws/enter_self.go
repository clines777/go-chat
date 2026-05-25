package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
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
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	ws.LeaveAllGroups(ctx.Client.ConnID)

	sess.Scene = protocol.SceneSelfInfo
	sess.InGroupID = 0
	if err := session.Set(ctx.Client.ConnID, sess); err != nil {
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	respData, err := json.Marshal(&protocol.UserSelfResp{
		Id:         u.ID,
		Username:   u.ExtUsername,
		Nickname:   u.Nickname,
		CreateTime: u.CreateTime,
		Code:       u.Code,
		UserLevel:  u.UserLevel,
	})
	if err != nil {
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	return &protocol.Payload{MsgType: protocol.EnterSelfOk, Data: respData, Meta: ctx.Payload.Meta}
}
