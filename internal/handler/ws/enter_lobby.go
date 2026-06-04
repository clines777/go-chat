package handler

import (
	"encoding/json"
	"log"

	"gochat/internal/group"
	"gochat/internal/protocol"
	"gochat/internal/session"
	"gochat/internal/ws"
)

// EnterLobby 用戶進入大廳場景
func EnterLobby(ctx *ws.Ctx) *protocol.Payload {
	var req protocol.EnterLobbyReq
	if ctx.Payload.Data != nil {
		if err := json.Unmarshal(ctx.Payload.Data, &req); err != nil {
			log.Printf("[EnterLobby] bad request: %v", err)
			return protocol.NewErrPayload(protocol.ErrInvalidParam, "invalid request", ctx.Payload)
		}
	}

	sess := session.Get(ctx.Client.UserId, ctx.Client.ConnID)

	groups, err := group.GetLobbyGroups(req.Page)
	if err != nil {
		log.Printf("[EnterLobby] GetLobbyGroups page=%d error: %v", req.Page, err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	ws.ResetGroupScene(ctx.Client.ConnID)

	sess.Scene = protocol.SceneLobby
	sess.InGroupID = 0
	if err := session.Set(ctx.Client.UserId, ctx.Client.ConnID, sess); err != nil {
		log.Printf("[EnterLobby] session.Set error: %v", err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	respData, err := json.Marshal(&protocol.EnterLobbyResp{
		Page:   req.Page,
		Groups: groups,
	})
	if err != nil {
		log.Printf("[EnterLobby] marshal error: %v", err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	return &protocol.Payload{MsgType: protocol.EnterLobbyOk, Data: respData, Meta: ctx.Payload.Meta}
}
