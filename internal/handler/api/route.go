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
		groupGroup.GET("/lobby", GetLobbyGroup)
		groupGroup.GET("/my", GetMyGroup)
		groupGroup.GET("/info", GetGroupInfo)
		groupGroup.POST("/join", JoinGroup)
		groupGroup.POST("/quit", QuitGroup)
		groupGroup.POST("/ban-user", BanUser)
		groupGroup.POST("/unban-user", UnbanUser)
		groupGroup.POST("/kick-user", KickUser)
	}

	chatGroup := apiRoute.Group("/chat")
	{
		chatGroup.GET("/sync", SyncChat)
	}

	adminRoute := e.Group("/admin")

	adminGroupGroup := adminRoute.Group("/group")
	{
		adminGroupGroup.POST("/create", CreateGroup)
		adminGroupGroup.POST("/update", UpdateGroup)
		adminGroupGroup.POST("/invite", InviteUser)
		adminGroupGroup.GET("/user", GetUserOfGroup)
	}

	adminUserGroup := adminRoute.Group("/user")
	{
		adminUserGroup.GET("/list", GetUserList)
		adminUserGroup.GET("/group", GetGroupOfUser)
		adminUserGroup.POST("/forbid", ForbidUser)
	}
}
