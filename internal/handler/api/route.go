package api

import (
	"github.com/gin-gonic/gin"
	"gochat/internal/protocol"
)

func RegisterRoutes(e *gin.Engine) {

	apiRoute := e.Group("/api")
	{
		//ping
		e.GET("/ping", func(c *gin.Context) {
			resp := &protocol.ApiResponse{Code: protocol.ErrorNone, Message: "OK", Data: nil}

			c.JSON(200, resp.Get())
		})

		//驗證
		apiRoute.Group("/auth")
		{
			//取登入Token
			apiRoute.POST("/get-login-token", GetLoginToken)
		}

		//用戶
		apiRoute.Group("/user")
		{
			//取用戶基本資訊.
			apiRoute.GET("/info", GetUserInfo)
		}

		//頭像
		apiRoute.Group("/avatar")
		{
			//頭像列表
			apiRoute.GET("/list", GetAvatarList)
			//設置頭像
			apiRoute.POST("/set", SetAvatar)
		}

		//群組
		apiRoute.Group("/group")
		{
			//大廳群組列表
			apiRoute.GET("/lobby", GetLobbyGroup)
			//我的群組列表
			apiRoute.GET("/my", GetMyGroup)
			//群組基本資訊.
			apiRoute.GET("/info", GetGroupInfo)
			//加入群組
			apiRoute.POST("/join", JoinGroup)
			//退出群組
			apiRoute.POST("/quit", QuitGroup)
			//用戶禁言
			apiRoute.POST("/ban-user", BanUser)
			//用戶解禁
			apiRoute.POST("/unban-user", UnbanUser)
			//踢用戶出群
			apiRoute.POST("/kick-user", KickUser)
		}

		//聊天訊息
		apiRoute.Group("/chat")
		{
			//同步歷史訊息
			apiRoute.GET("/sync", SyncChat)
		}

		//後台接口
		adminRoute := e.Group("/admin")
		{
			//群組相關
			adminRoute.Group("/group")
			{
				//創建群組
				adminRoute.POST("/create", CreateGroup)
				//更新群組
				adminRoute.POST("/update", UpdateGroup)
				//邀請用戶進群
				adminRoute.POST("/invite", InviteUser)
				//群組用戶
				adminRoute.GET("/user", GetUserOfGroup)
			}

			//用戶相關
			adminRoute.Group("/user")
			{
				//用戶列表
				adminRoute.GET("/list", GetUserList)
				//用戶的群組列表
				adminRoute.GET("/group", GetGroupOfUser)
				//禁用
				adminRoute.POST("/forbid", ForbidUser)
			}
		}
	}
}
