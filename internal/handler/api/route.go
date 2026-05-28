package api

import (
	"github.com/gin-gonic/gin"
	"gochat/internal/middleware"
	"gochat/internal/protocol"
)

func RegisterRoutes(e *gin.Engine) {
	e.Static("/static", "./static")
	e.StaticFile("/", "./web/index.html")

	e.GET("/ping", func(c *gin.Context) {
		c.JSON(200, (&protocol.ApiResponse{Code: protocol.ErrorNone, Message: "OK"}).H())
	})

	apiRoute := e.Group("/api")

	authGroup := apiRoute.Group("/auth")
	{
		authGroup.POST("/get-login-token", GetLoginToken)
	}

	protected := apiRoute.Group("", middleware.Auth())

	userGroup := protected.Group("/user")
	{
		userGroup.GET("/self", GetUserSelfInfo)
	}

	avatarGroup := protected.Group("/avatar")
	{
		avatarGroup.GET("/list", GetAvatars)
		avatarGroup.POST("/set", SetAvatar)
	}

	groupGroup := protected.Group("/group")
	{
		groupGroup.POST("/create", CreateGroup)
		groupGroup.GET("/info", GetGroupInfo)
		groupGroup.POST("/join", JoinGroup)
		groupGroup.POST("/cover/upload", UploadGroupCover)
	}
}
