package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gochat/internal/group"
	"gochat/internal/protocol"
	"gochat/internal/user"
)

// GetGroupInfo - 取得群組資訊頁面展示數據
func GetGroupInfo(c *gin.Context) {
	var req protocol.GetGroupInfoReq
	if err := c.ShouldBindQuery(&req); err != nil || req.GroupID == 0 {
		c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrWrongParam, "參數錯誤", nil).H())
		return
	}

	g, err := group.FindByID(req.GroupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorNotFound, "群組不存在", nil).H())
		} else {
			log.Printf("[GetGroupInfo] FindByID group=%d error: %v", req.GroupID, err)
			c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorUnknown, "系統錯誤", nil).H())
		}
		return
	}

	count, err := group.GetMemberCount(g.ID)
	if err != nil {
		log.Printf("[GetGroupInfo] GetMemberCount group=%d error: %v", g.ID, err)
		c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorUnknown, "系統錯誤", nil).H())
		return
	}

	c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorNone, "OK", &protocol.GroupInfoResp{
		Title:         g.Title,
		UserTotal:     int32(count),
		Bulletin:      g.Bulletin,
		OwnerUsername: g.OwnerUserName,
		OwnerUserID:   int64(g.OwnerUserID),
		Code:          g.Code,
		Remark:        g.Remark,
	}).H())
}

// JoinGroup - 用戶加入群組
func JoinGroup(c *gin.Context) {
	var req protocol.JoinGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrWrongParam, "參數錯誤", nil).H())
		return
	}

	userID := c.MustGet("user_id").(int64)

	u, err := user.FindByID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrUserNotFound, "用戶不存在", nil).H())
		} else {
			log.Printf("[JoinGroup] FindByID user=%d error: %v", userID, err)
			c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorUnknown, "系統錯誤", nil).H())
		}
		return
	}

	g, err := group.FindByID(req.GroupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorNotFound, "群組不存在", nil).H())
		} else {
			log.Printf("[JoinGroup] FindByID group=%d error: %v", req.GroupID, err)
			c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorUnknown, "系統錯誤", nil).H())
		}
		return
	}

	if err := group.Join(u, g); err != nil {
		switch {
		case errors.Is(err, group.ErrAlreadyMember):
			c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrInvalidParam, "已是群組成員", nil).H())
		case errors.Is(err, group.ErrGroupFull):
			c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrInvalidParam, "群組已達人數上限", nil).H())
		case errors.Is(err, group.ErrGroupNotOpen):
			c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrInvalidParam, "群組未開放加入", nil).H())
		case errors.Is(err, group.ErrGroupDismissed):
			c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrInvalidParam, "群組已解散", nil).H())
		default:
			log.Printf("[JoinGroup] Join user=%d group=%d error: %v", userID, req.GroupID, err)
			c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorUnknown, "系統錯誤", nil).H())
		}
		return
	}

	c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorNone, "OK", nil).H())
}

// CreateGroup - 創建群組
func CreateGroup(c *gin.Context) {
	var req protocol.CreateGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrWrongParam, "參數錯誤", nil).H())
		return
	}

	userID := c.MustGet("user_id").(int64)

	owner, err := user.FindByID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrUserNotFound, "用戶不存在", nil).H())
		} else {
			log.Printf("[CreateGroup] FindByID user=%d error: %v", userID, err)
			c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorUnknown, "系統錯誤", nil).H())
		}
		return
	}

	groupID, code, err := group.Create(&req, owner)
	if err != nil {
		log.Printf("[CreateGroup] Create user=%d title=%q error: %v", userID, req.Title, err)
		c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorUnknown, "建立群組失敗", nil).H())
		return
	}

	c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorNone, "OK", map[string]any{
		"group_id": groupID,
		"code":     code,
	}).H())
}
