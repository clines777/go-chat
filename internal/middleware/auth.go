package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gochat/internal/infra/redis"
	"gochat/internal/protocol"
)

// Auth - 驗證Bearer token並取得 user id後寫入Context供後續使用.
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		token, ok := strings.CutPrefix(header, "Bearer ")
		if !ok || token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, protocol.NewApiResponse(protocol.ErrUnauthorized, "unauthorized", nil).H())
			return
		}

		var info protocol.ApiTokenInfo
		found, err := redis.GetRedis().GetJSON(protocol.ApiTokenKey(token), &info)
		if err != nil || !found || info.UserID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, protocol.NewApiResponse(protocol.ErrUnauthorized, "unauthorized", nil).H())
			return
		}

		c.Set("user_id", info.UserID)
		c.Next()
	}
}
