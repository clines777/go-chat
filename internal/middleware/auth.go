package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gochat/internal/infra/redis"
	"gochat/internal/protocol"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		token, ok := strings.CutPrefix(header, "Bearer ")
		if !ok || token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, protocol.NewApiResponse(protocol.ErrUnauthorized, "unauthorized", nil).H())
			return
		}

		var payload protocol.ApiTokenPayload
		found, err := redis.GetRedis().GetJSON(protocol.ApiTokenKey(token), &payload)
		if err != nil || !found || payload.UserID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, protocol.NewApiResponse(protocol.ErrUnauthorized, "unauthorized", nil).H())
			return
		}

		c.Set("user_id", payload.UserID)
		c.Next()
	}
}
