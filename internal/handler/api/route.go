package api

import (
	"github.com/gin-gonic/gin"
	"gochat/internal/middleware"
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

	protected := apiRoute.Group("", middleware.Auth())

	userGroup := protected.Group("/user")
	{
		userGroup.GET("/info", GetUserInfo)
	}

	avatarGroup := protected.Group("/avatar")
	{
		avatarGroup.GET("/list", GetAvatarList)
		avatarGroup.POST("/set", SetAvatar)
	}

	groupGroup := protected.Group("/group")
	{
		groupGroup.POST("/create", CreateGroup)
		groupGroup.GET("/info", GetGroupInfo)
		groupGroup.POST("/join", JoinGroup)
	}
}
