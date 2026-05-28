package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gochat/internal/infra/db"
	"gochat/internal/model"
	"gochat/internal/protocol"
	"gochat/internal/user"
)

type avatarResp struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

func GetAvatars(c *gin.Context) {
	avatars, err := getAvatarList()
	if err != nil {
		log.Printf("[GetAvatars] getAvatarList error: %v", err)
		c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorUnknown, "系統錯誤", []avatarResp{}).H())
		return
	}

	resp := make([]avatarResp, len(avatars))
	for i, a := range avatars {
		resp[i] = avatarResp{ID: a.ID, URL: "/static/avatars/" + a.Filename}
	}

	c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorNone, "OK", resp).H())
}

func SetAvatar(c *gin.Context) {
	var req protocol.SetAvatarReq
	if err := c.ShouldBindJSON(&req); err != nil || req.AvatarId == 0 {
		c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrWrongParam, "參數錯誤", nil).H())
		return
	}

	userID := c.MustGet("user_id").(int)

	if err := user.UpdateAvatar(userID, req.AvatarId); err != nil {
		log.Printf("[SetAvatar] UpdateAvatar user=%d avatar=%d error: %v", userID, req.AvatarId, err)
		c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorUnknown, "系統錯誤", nil).H())
		return
	}

	c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorNone, "OK", nil).H())
}

func getAvatarList() ([]model.Avatar, error) {
	d, err := db.GetDBConn()
	if err != nil {
		return nil, err
	}

	rows := make([]model.Avatar, 0)
	querySql, args, _ := d.Builder.
		Select("id", "filename").
		From("avatar").
		ToSql()

	err = d.DB.Select(&rows, querySql, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return rows, nil
		}
		return nil, err
	}

	return rows, nil
}
