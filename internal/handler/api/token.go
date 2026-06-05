package api

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gochat/internal/infra/redis"
	"gochat/internal/protocol"
)

// GetLoginToken 取得登入token, 此token會暫放redis, 並用於websocket登入, 防止刷號跟刷連線.(先省掉密碼了, )
func GetLoginToken(c *gin.Context) {
	var req protocol.GetTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrWrongParam, "參數錯誤", nil).H())
		return
	}

	if req.Username == "" {
		c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrWrongParam, "username 不可為空", nil).H())
		return
	}

	token, err := genLoginToken()
	if err != nil {
		log.Printf("[GetLoginToken] genLoginToken user=%s error: %v", req.Username, err)
		c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorUnknown, "系統錯誤", nil).H())
		return
	}
	if err := redis.GetRedis().SetJSON(protocol.LoginTokenKey(token), req, 10*time.Second); err != nil {
		log.Printf("[GetLoginToken] SetJSON user=%s error: %v", req.Username, err)
		c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorUnknown, "系統錯誤", nil).H())
		return
	}

	c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorNone, "OK", map[string]any{
		"token": token,
	}).H())
}

// genLoginToken 由於後面版本部分參數拿掉, 直接改用random
func genLoginToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return strings.ToUpper(hex.EncodeToString(b)), nil
}
