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
	"gochat/internal/user"
	"gochat/internal/ws"
)

// JoinGroup - 用戶加入群組成員
func JoinGroup(ctx *ws.Ctx) *protocol.Payload {
	var req protocol.JoinGroupReq
	if err := json.Unmarshal(ctx.Payload.Data, &req); err != nil || req.GroupID <= 0 {
		return protocol.NewErrPayload(protocol.ErrInvalidParam, "invalid request", ctx.Payload)
	}

	sess := session.Get(ctx.Client.UserId, ctx.Client.ConnID)

	u, err := user.FindByID(sess.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return protocol.NewErrPayload(protocol.ErrUserNotFound, "用戶不存在", ctx.Payload)
		}
		log.Printf("[JoinGroup] user.FindByID user=%d error: %v", sess.UserID, err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	g, err := group.FindByID(req.GroupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return protocol.NewErrPayload(protocol.ErrorNotFound, "群組不存在", ctx.Payload)
		}
		log.Printf("[JoinGroup] group.FindByID group=%d error: %v", req.GroupID, err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	if err := group.Join(u, g); err != nil {
		switch {
		case errors.Is(err, group.ErrAlreadyMember):
			return protocol.NewErrPayload(protocol.ErrInvalidParam, "已是群組成員", ctx.Payload)
		case errors.Is(err, group.ErrGroupFull):
			return protocol.NewErrPayload(protocol.ErrInvalidParam, "群組已達人數上限", ctx.Payload)
		case errors.Is(err, group.ErrGroupNotOpen):
			return protocol.NewErrPayload(protocol.ErrInvalidParam, "群組未開放加入", ctx.Payload)
		case errors.Is(err, group.ErrGroupDismissed):
			return protocol.NewErrPayload(protocol.ErrInvalidParam, "群組已解散", ctx.Payload)
		default:
			log.Printf("[JoinGroup] Join user=%d group=%d error: %v", sess.UserID, req.GroupID, err)
			return protocol.NewErrPayload(protocol.ErrInternalError, "加入失敗", ctx.Payload)
		}
	}

	// 加入成功後自動進入群組
	count, err := group.GetMemberCount(req.GroupID)
	if err != nil {
		log.Printf("[JoinGroup] GetMemberCount group=%d error: %v", req.GroupID, err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	ws.JoinGroup(ctx.Client.ConnID, req.GroupID, ctx.Client)

	sess.InGroupID = req.GroupID
	sess.Scene = protocol.SceneInGroup
	if err := session.Set(ctx.Client.UserId, ctx.Client.ConnID, sess); err != nil {
		log.Printf("[JoinGroup] session.Set error: %v", err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	chats, hasMore, err := chat.GetRecentChats(req.GroupID, RecentChatCount)
	if err != nil {
		log.Printf("[JoinGroup] GetRecentChats group=%d error: %v", req.GroupID, err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	pinChat := FindPinChat(req.GroupID, g.PinChatID)
	respData, err := json.Marshal(&protocol.EnterGroupResp{
		Title:          g.Title,
		GroupId:        req.GroupID,
		GroupUserCount: count,
		PinChat:        pinChat,
		SelfBanned:     false,
		Chats:          chats,
		HasMore:        hasMore,
	})
	if err != nil {
		log.Printf("[JoinGroup] marshal error: %v", err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	log.Printf("[JoinGroup] OK user=%d joined group=%d", sess.UserID, req.GroupID)
	return &protocol.Payload{MsgType: protocol.JoinGroupOK, Data: respData, Meta: ctx.Payload.Meta}
}
