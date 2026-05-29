package handler

import (
	"encoding/json"
	"log"

	"gochat/internal/group"
	"gochat/internal/model"
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
		log.Printf("[Resume] token not exists or expired")
		return protocol.NewErrPayload(protocol.ErrUnauthorized, "invalid or expired resume token", ctx.Payload)
	}

	exists := user.Exists(tokenPayload.UserID)
	if !exists {
		log.Printf("[Resume] user not exists")
		return protocol.NewErrPayload(protocol.ErrUserNotFound, "user not exists", ctx.Payload)
	}

	ctx.Client.UserId = tokenPayload.UserID
	ws.Register(ctx.Client)

	userGroups, err := group.GetMyGroups(tokenPayload.UserID)
	if err != nil {
		log.Printf("[Resume] GetMyGroups user=%d error: %v", tokenPayload.UserID, err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "fetch groups error", ctx.Payload)
	}

	apiToken, err := user.GenerateApiToken(tokenPayload.UserID)
	if err != nil {
		log.Printf("[Resume] GenerateApiToken error: %v", err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}
	ctx.Client.ApiToken = apiToken

	user.DeleteResumeToken(req.Token)
	newResumeToken, err := user.GenerateResumeToken(&model.User{ID: tokenPayload.UserID, Username: tokenPayload.Username})
	if err != nil {
		log.Printf("[Resume] GenerateResumeToken error: %v", err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	//意外斷線後跳回我的群組
	sess := &session.Session{
		ConnID:   ctx.Client.ConnID,
		UserID:   tokenPayload.UserID,
		Username: tokenPayload.Username,
		Scene:    protocol.SceneMyGroup,
		ApiToken: apiToken,
	}
	if err := session.Set(tokenPayload.UserID, ctx.Client.ConnID, sess); err != nil {
		log.Printf("[Resume] session.Set error: %v", err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	respData, err := json.Marshal(&protocol.LoginResp{
		UserID:      tokenPayload.UserID,
		Username:    tokenPayload.Username,
		ApiToken:    apiToken,
		ResumeToken: newResumeToken,
		UserGroups:  userGroups,
	})
	if err != nil {
		log.Printf("[Resume] marshal error: %v", err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	log.Printf("[Resume] OK user=%d (%s)", tokenPayload.UserID, tokenPayload.Username)
	return &protocol.Payload{MsgType: protocol.ResumeOK, Data: respData, Meta: ctx.Payload.Meta}
}
