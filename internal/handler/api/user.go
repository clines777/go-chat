package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gochat/internal/protocol"
	"gochat/internal/user"
)

func GetUserSelfInfo(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)

	u, err := user.FindByID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrUserNotFound, "用戶不存在", nil).H())
		} else {
			c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorUnknown, "系統錯誤", nil).H())
		}
		return
	}

	avatarURL := ""
	if u.AvatarFilename != "" {
		avatarURL = "/static/avatars/" + u.AvatarFilename
	}

	c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorNone, "OK", &protocol.UserSelfResp{
		Id:         u.ID,
		Username:   u.ExtUsername,
		AvatarURL:  avatarURL,
		CreateTime: u.CreateTime,
		Code:       u.Code,
		UserLevel:  u.UserLevel,
	}).H())
}

func ForbidUser(c *gin.Context) {

}
