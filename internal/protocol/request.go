package protocol

import _ "github.com/go-playground/validator/v10"

type GetTokenReq struct {
	SiteBid  string `json:"site_bid"`
	Username string `json:"username"`
	MemberId int32  `json:"member_id"`
}

type GetUserInfoReq struct {
	UserId int32 `form:"user_id"`
}

type ApiToken struct {
	ApiToken string `json:"api_token"`
}
