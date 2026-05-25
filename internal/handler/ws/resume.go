package handler

import (
	"encoding/json"
	"gochat/internal/group"
	"gochat/internal/protocol"
	"gochat/internal/session"
	"gochat/internal/user"
	"gochat/internal/ws"
)

func Resume(ctx *ws.Ctx) *protocol.Payload {
	var req protocol.ResumeReq
	if err := json.Unmarshal(ctx.Payload.Data, &req); err != nil || req.Token == "" {
		return protocol.NewErrPayload(protocol.ErrInvalidParam, "invalid request", ctx.Payload)
	}

	tokenPayload, err := user.GetResumeToken(req.Token)
	if err != nil {
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}
	if tokenPayload == nil {
		return protocol.NewErrPayload(protocol.ErrUnauthorized, "invalid or expired token", ctx.Payload)
	}

	u, err := user.FindByID(tokenPayload.UserID)
	if err != nil {
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
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}
	ws.Register(ctx.Client)

	userGroups, err := group.GetMyGroups(tokenPayload.UserID, tokenPayload.SiteBid)
	if err != nil {
		return protocol.NewErrPayload(protocol.ErrInternalError, "fetch groups error", ctx.Payload)
	}

	apiToken, err := user.GenerateApiToken(tokenPayload.UserID, tokenPayload.SiteBid)
	if err != nil {
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
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	return &protocol.Payload{MsgType: protocol.ResumeOK, Data: respData, Meta: ctx.Payload.Meta}
}
