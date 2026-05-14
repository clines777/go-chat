package api

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gochat/internal/infra/redis"
	"gochat/internal/protocol"
)

const salt = "CsXGM2ArECWqoLT0BKst"

// GetLoginToken - 包網端登入用API, 傳入帳號資訊順帶生成websocket登入用Token放redis, 先省略加密
func GetLoginToken(c *gin.Context) {
	var req protocol.GetTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrWrongParam, "參數錯誤", nil).Get())
		return
	}

	if req.SiteBid == "" || req.Username == "" || req.MemberId == 0 {
		c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrWrongParam, "site_bid, username, member_id 不可為空", nil).Get())
		return
	}

	token := genLoginToken(&req)
	if err := redis.GetRedis().SetJSON(protocol.LoginTokenKey(token), req, 10*time.Second); err != nil {
		c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorUnknown, "系統錯誤", nil).Get())
		return
	}

	c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorNone, "OK", map[string]any{
		"token": token,
	}).Get())
}

func genLoginToken(req *protocol.GetTokenReq) string {
	info := fmt.Sprintf("%v%v%v%v%v", req.SiteBid, req.MemberId, req.Username, time.Now(), salt)
	hash := md5.Sum([]byte(info))
	return strings.ToUpper(req.SiteBid + hex.EncodeToString(hash[:]))
}
