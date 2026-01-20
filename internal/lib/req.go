package lib

type GetLoginTokenReq struct {
	MemberId     int32  `json:"member_id"`
	Username     string `json:"username"`
	SiteBid      string `json:"site_bid"`
	PlatformType int8   `json:"platform_type"`
}

type GetUserInfoReq struct {
	UserId int32 `form:"user_id"`
}

type ApiToken struct {
	ApiToken string `json:"api_token"`
}
