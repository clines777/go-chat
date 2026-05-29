package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"

	"gochat/internal/chat"
	"gochat/internal/group"
	"gochat/internal/protocol"
	"gochat/internal/session"
	"gochat/internal/ws"
)

const RecentChatCount = 30

func EnterGroup(ctx *ws.Ctx) *protocol.Payload {
	var req protocol.EnterGroupReq
	if err := json.Unmarshal(ctx.Payload.Data, &req); err != nil || req.GroupID == 0 {
		log.Printf("[EnterGroup] bad request: %v", err)
		return protocol.NewErrPayload(protocol.ErrInvalidParam, "invalid request", ctx.Payload)
	}

	sess := session.Get(ctx.Client.UserId, ctx.Client.ConnID)

	_, err := group.GetMembership(sess.UserID, req.GroupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return protocol.NewErrPayload(protocol.ErrNotMember, "not a member", ctx.Payload)
		}
		log.Printf("[EnterGroup] GetMembership user=%d group=%d error: %v", sess.UserID, req.GroupID, err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	g, err := group.FindByID(req.GroupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return protocol.NewErrPayload(protocol.ErrorNotFound, "group not found", ctx.Payload)
		}
		log.Printf("[EnterGroup] FindByID group=%d error: %v", req.GroupID, err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	count, err := group.GetMemberCount(req.GroupID)
	if err != nil {
		log.Printf("[EnterGroup] GetMemberCount group=%d error: %v", req.GroupID, err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	ws.JoinGroup(ctx.Client.ConnID, req.GroupID, ctx.Client)

	sess.InGroupID = req.GroupID
	sess.Scene = protocol.SceneInGroup
	if err := session.Set(ctx.Client.UserId, ctx.Client.ConnID, sess); err != nil {
		log.Printf("[EnterGroup] session.Set error: %v", err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	chats, err := chat.GetRecentChats(req.GroupID, RecentChatCount)
	if err != nil {
		log.Printf("[EnterGroup] GetRecentChats group=%d error: %v", req.GroupID, err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	respData, err := json.Marshal(&protocol.EnterGroupResp{
		Title:          g.Title,
		GroupId:        req.GroupID,
		GroupUserCount: count,
		Chats:          chats,
	})
	if err != nil {
		log.Printf("[EnterGroup] marshal error: %v", err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	return &protocol.Payload{MsgType: protocol.EnterGroupOk, Data: respData, Meta: ctx.Payload.Meta}
}
