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

func SendChat(ctx *ws.Ctx) *protocol.Payload {
	var req protocol.SendChatReq
	if err := json.Unmarshal(ctx.Payload.Data, &req); err != nil || req.GroupId == 0 || req.Content == "" {
		log.Printf("[SendChat] bad request: %v", err)
		return protocol.NewErrPayload(protocol.ErrInvalidParam, "invalid request", ctx.Payload)
	}

	sess := session.Get(ctx.Client.UserId, ctx.Client.ConnID)
	if sess.Scene != protocol.SceneInGroup || sess.InGroupID != req.GroupId {
		return protocol.NewErrPayload(protocol.ErrInvalidParam, "not in group", ctx.Payload)
	}

	g, err := group.FindByID(req.GroupId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return protocol.NewErrPayload(protocol.ErrorNotFound, "group not found", ctx.Payload)
		}
		log.Printf("[SendChat] FindByID group=%d error: %v", req.GroupId, err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}
	if g.IsDismiss {
		return protocol.NewErrPayload(protocol.ErrInvalidParam, "group is dismissed", ctx.Payload)
	}
	record, err := chat.SaveRecord(req.GroupId, sess.UserID, req.Content)
	if err != nil {
		log.Printf("[SendChat] SaveRecord group=%d user=%d error: %v", req.GroupId, sess.UserID, err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	avatarURL := ""
	if u, err := user.FindByID(sess.UserID); err == nil && u.AvatarFilename != "" {
		avatarURL = "/static/avatars/" + u.AvatarFilename
	}

	event := &protocol.CastChatEvent{
		Id:         record.ID,
		GroupId:    req.GroupId,
		UserId:     sess.UserID,
		Username:   sess.Username,
		AvatarURL:  avatarURL,
		Content:    req.Content,
		CreateTime: record.CreateTime,
	}
	if err := chat.Publish(event); err != nil {
		log.Printf("[SendChat] Publish group=%d error: %v", req.GroupId, err)
	}

	respData, err := json.Marshal(&protocol.SendChatResp{
		Id:         record.ID,
		GroupId:    req.GroupId,
		CreateTime: record.CreateTime,
	})
	if err != nil {
		log.Printf("[SendChat] marshal error: %v", err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	return &protocol.Payload{MsgType: protocol.SendChatOk, Data: respData, Meta: ctx.Payload.Meta}
}
