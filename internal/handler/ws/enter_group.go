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
	if err := json.Unmarshal(ctx.Payload.Data, &req); err != nil || req.GroupId == 0 {
		log.Printf("[EnterGroup] bad request: %v", err)
		return protocol.NewErrPayload(protocol.ErrInvalidParam, "invalid request", ctx.Payload)
	}

	sess := session.Get(ctx.Client.ConnID)

	membership, err := group.GetMembership(sess.UserID, int64(req.GroupId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return protocol.NewErrPayload(protocol.ErrNotMember, "not a member", ctx.Payload)
		}
		log.Printf("[EnterGroup] GetMembership user=%d group=%d error: %v", sess.UserID, req.GroupId, err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}
	if membership.IsBan {
		return protocol.NewErrPayload(protocol.ErrUserBanned, "user is banned", ctx.Payload)
	}

	g, err := group.FindByID(int64(req.GroupId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return protocol.NewErrPayload(protocol.ErrorNotFound, "group not found", ctx.Payload)
		}
		log.Printf("[EnterGroup] FindByID group=%d error: %v", req.GroupId, err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	count, err := group.GetMemberCount(int64(req.GroupId))
	if err != nil {
		log.Printf("[EnterGroup] GetMemberCount group=%d error: %v", req.GroupId, err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	ws.JoinGroup(ctx.Client.ConnID, req.GroupId, ctx.Client)

	sess.InGroupID = req.GroupId
	sess.Scene = protocol.SceneInGroup
	if err := session.Set(ctx.Client.ConnID, sess); err != nil {
		log.Printf("[EnterGroup] session.Set error: %v", err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	chats, err := chat.GetRecentChats(req.GroupId, RecentChatCount)
	if err != nil {
		log.Printf("[EnterGroup] GetRecentChats group=%d error: %v", req.GroupId, err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	respData, err := json.Marshal(&protocol.EnterGroupResp{
		Title:          g.Title,
		GroupId:        req.GroupId,
		GroupUserCount: int32(count),
		Chats:          chats,
	})
	if err != nil {
		log.Printf("[EnterGroup] marshal error: %v", err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	return &protocol.Payload{MsgType: protocol.EnterGroupOk, Data: respData, Meta: ctx.Payload.Meta}
}
