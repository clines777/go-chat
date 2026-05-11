package api

import (
	"github.com/gin-gonic/gin"
	"gochat/internal/protocol"
)

func RegisterRoutes(e *gin.Engine) {
	e.GET("/ping", func(c *gin.Context) {
		c.JSON(200, (&protocol.ApiResponse{Code: protocol.ErrorNone, Message: "OK"}).Get())
	})

	apiRoute := e.Group("/api")

	authGroup := apiRoute.Group("/auth")
	{
		authGroup.POST("/get-login-token", GetLoginToken)
	}

	userGroup := apiRoute.Group("/user")
	{
		userGroup.GET("/info", GetUserInfo)
	}

	avatarGroup := apiRoute.Group("/avatar")
	{
		avatarGroup.GET("/list", GetAvatarList)
		avatarGroup.POST("/set", SetAvatar)
	}

	groupGroup := apiRoute.Group("/group")
	{
		groupGroup.GET("/info", GetGroupInfo)
		groupGroup.POST("/join", JoinGroup)
	}
}
