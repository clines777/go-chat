package route

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) error {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "pong"})
	})

}
