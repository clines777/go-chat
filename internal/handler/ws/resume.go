package handler

import (
	"encoding/json"
	"log"

	"gochat/internal/group"
	"gochat/internal/protocol"
	"gochat/internal/session"
	"gochat/internal/user"
	"gochat/internal/ws"
)

func Resume(ctx *ws.Ctx) *protocol.Payload {
	var req protocol.ResumeReq
	if err := json.Unmarshal(ctx.Payload.Data, &req); err != nil || req.Token == "" {
		log.Printf("[Resume] bad request: %v", err)
		return protocol.NewErrPayload(protocol.ErrInvalidParam, "invalid request", ctx.Payload)
	}

	tokenPayload, err := user.GetResumeToken(req.Token)
	if err != nil {
		log.Printf("[Resume] GetResumeToken error: %v", err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}
	if tokenPayload == nil {
		log.Printf("[Resume] token not found or expired")
		return protocol.NewErrPayload(protocol.ErrUnauthorized, "invalid or expired token", ctx.Payload)
	}

	u, err := user.FindByID(tokenPayload.UserID)
	if err != nil {
		log.Printf("[Resume] FindByID user=%d error: %v", tokenPayload.UserID, err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	sess := &session.Session{
		ConnID:    ctx.Client.ConnID,
		UserID:    tokenPayload.UserID,
		SiteBid:   tokenPayload.SiteBid,
		Username:  tokenPayload.Username,
		UserLevel: u.UserLevel,
		Scene:     protocol.SceneMyGroup,
	}
	if err := session.Set(ctx.Client.ConnID, sess); err != nil {
		log.Printf("[Resume] session.Set error: %v", err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}
	ws.Register(ctx.Client)

	userGroups, err := group.GetMyGroups(tokenPayload.UserID, tokenPayload.SiteBid)
	if err != nil {
		log.Printf("[Resume] GetMyGroups user=%d error: %v", tokenPayload.UserID, err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "fetch groups error", ctx.Payload)
	}

	apiToken, err := user.GenerateApiToken(tokenPayload.UserID, tokenPayload.SiteBid)
	if err != nil {
		log.Printf("[Resume] GenerateApiToken error: %v", err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	user.RefreshResumeToken(req.Token)

	respData, err := json.Marshal(&protocol.LoginResp{
		UserID:      tokenPayload.UserID,
		Username:    tokenPayload.Username,
		ApiToken:    apiToken,
		ResumeToken: req.Token,
		UserGroups:  userGroups,
	})
	if err != nil {
		log.Printf("[Resume] marshal error: %v", err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	log.Printf("[Resume] OK user=%d (%s)", tokenPayload.UserID, tokenPayload.Username)
	return &protocol.Payload{MsgType: protocol.ResumeOK, Data: respData, Meta: ctx.Payload.Meta}
}
