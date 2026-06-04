package api

import "github.com/gin-gonic/gin"

// getContextUID 取出auth middleware中解析並寫入的 user_id。
func getContextUID(c *gin.Context) (int, bool) {
	v, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	id, ok := v.(int)
	if !ok {
		return 0, false
	}
	return id, true
}
