package route

import (
	"github.com/gin-gonic/gin"
	"gochat/internal/code"
	"gochat/internal/handler/api"
	"gochat/internal/lib"
)

func RegisterRoutes(e *gin.Engine) {

	apiRoute := e.Group("/api")
	{
		//ping
		e.GET("/ping", func(c *gin.Context) {
			resp := &lib.ApiResponse{Code: code.ErrorNone, Message: "OK", Data: nil}

			c.JSON(200, resp.Get())
		})

		//驗證
		apiRoute.Group("/auth")
		{
			//取登入Token
			apiRoute.POST("/get-login-token", api.GetLoginToken)
		}

		//用戶
		apiRoute.Group("/user")
		{
			//取用戶基本資訊.
			apiRoute.GET("/info", api.GetUserInfo)
		}

		//頭像
		apiRoute.Group("/avatar")
		{
			//頭像列表
			apiRoute.GET("/list", api.GetAvatarList)
			//設置頭像
			apiRoute.POST("/set", api.SetAvatar)
		}

		//群組
		apiRoute.Group("/group")
		{
			//大廳群組列表
			apiRoute.GET("/lobby", api.GetLobbyGroup)
			//我的群組列表
			apiRoute.GET("/my", api.GetMyGroup)
			//群組基本資訊.
			apiRoute.GET("/info", api.GetGroupInfo)
			//加入群組
			apiRoute.POST("/join", api.JoinGroup)
			//退出群組
			apiRoute.POST("/quit", api.QuitGroup)
			//用戶禁言
			apiRoute.POST("/ban-user", api.BanUser)
			//用戶解禁
			apiRoute.POST("/unban-user", api.UnbanUser)
			//踢用戶出群
			apiRoute.POST("/kick-user", api.KickUser)
		}

		//聊天訊息
		apiRoute.Group("/chat")
		{
			//同步歷史訊息
			apiRoute.GET("/sync", api.SyncChat)
		}

		//後台接口
		adminRoute := e.Group("/admin")
		{
			//群組相關
			adminRoute.Group("/group")
			{
				//創建群組
				adminRoute.POST("/create", api.CreateGroup)
				//更新群組
				adminRoute.POST("/update", api.UpdateGroup)
				//邀請用戶進群
				adminRoute.POST("/invite", api.InviteUser)
				//群組用戶
				adminRoute.GET("/user", api.GetUserOfGroup)
			}

			//用戶相關
			adminRoute.Group("/user")
			{
				//用戶列表
				adminRoute.GET("/list", api.GetUserList)
				//用戶的群組列表
				adminRoute.GET("/group", api.GetGroupOfUser)
				//禁用
				adminRoute.POST("/forbid", api.ForbidUser)
			}
		}
	}
}
