package handler

import (
	"encoding/json"
	"log"

	"gochat/internal/group"
	"gochat/internal/protocol"
	"gochat/internal/user"
	"gochat/internal/ws"
)

func Login(ctx *ws.Ctx) *protocol.Payload {
	var req protocol.LoginReq
	if err := json.Unmarshal(ctx.Payload.Data, &req); err != nil || req.Token == "" {
		log.Printf("[Login] invalid token: err=%v ", err)
		return protocol.NewErrPayload(protocol.ErrWrongParam, "invalid request", ctx.Payload)
	}

	tokenInfo, err := user.GetLoginToken(req)
	if err != nil || tokenInfo == nil {
		log.Printf("[Login] invalid token: err=%v tokenFound=%v", err, tokenInfo != nil)
		return protocol.NewErrPayload(protocol.ErrUnauthorized, "invalid token", ctx.Payload)
	}

	u, err := user.Login(ctx, tokenInfo)
	if err != nil {
		log.Printf("[Login] user.Login error: %v", err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "user login error", ctx.Payload)
	}
	ws.Register(ctx.Client)

	userGroups, err := group.GetMyGroups(u.ID)
	if err != nil {
		log.Printf("[Login] GetMyGroups error: %v", err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "fetch groups error", ctx.Payload)
	}

	apiToken, err := user.GenerateApiToken(u.ID)
	if err != nil {
		log.Printf("[Login] GenerateApiToken error: %v", err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	resumeToken, err := user.GenerateResumeToken(u)
	if err != nil {
		log.Printf("[Login] GenerateResumeToken error: %v", err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	respData, err := json.Marshal(&protocol.LoginResp{
		UserID:      u.ID,
		Username:    u.Username,
		ApiToken:    apiToken,
		ResumeToken: resumeToken,
		UserGroups:  userGroups,
	})
	if err != nil {
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	log.Printf("[Login] OK user=%d (%s)", u.ID, u.Username)
	return &protocol.Payload{MsgType: protocol.LoginOK, Data: respData, Meta: ctx.Payload.Meta}
}
