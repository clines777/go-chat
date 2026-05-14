package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gochat/internal/group"
	"gochat/internal/protocol"
	"gochat/internal/user"
)

func GetGroupInfo(c *gin.Context) {

}

// JoinGroup - 用戶加入群組
func JoinGroup(c *gin.Context) {
	var req protocol.JoinGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrWrongParam, "參數錯誤", nil).Get())
		return
	}

	u, err := user.FindByID(req.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrUserNotFound, "用戶不存在", nil).Get())
		} else {
			c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorUnknown, "系統錯誤", nil).Get())
		}
		return
	}

	g, err := group.FindByID(req.GroupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorNotFound, "群組不存在", nil).Get())
		} else {
			c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorUnknown, "系統錯誤", nil).Get())
		}
		return
	}

	if err := group.Join(u, g); err != nil {
		switch {
		case errors.Is(err, group.ErrAlreadyMember):
			c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrInvalidParam, "已是群組成員", nil).Get())
		case errors.Is(err, group.ErrGroupFull):
			c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrInvalidParam, "群組已達人數上限", nil).Get())
		case errors.Is(err, group.ErrLevelInsufficient):
			c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrInvalidParam, "用戶等級不足", nil).Get())
		case errors.Is(err, group.ErrGroupNotOpen):
			c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrInvalidParam, "群組未開放加入", nil).Get())
		case errors.Is(err, group.ErrGroupDismissed):
			c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrInvalidParam, "群組已解散", nil).Get())
		case errors.Is(err, group.ErrSiteMismatch):
			c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrInvalidParam, "站台不符", nil).Get())
		default:
			c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorUnknown, "系統錯誤", nil).Get())
		}
		return
	}

	c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorNone, "OK", nil).Get())
}

// CreateGroup 創建群組
func CreateGroup(c *gin.Context) {
	var req protocol.CreateGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrWrongParam, "參數錯誤", nil).Get())
		return
	}

	owner, err := user.FindByID(req.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorNotFound, "用戶不存在", nil).Get())
		} else {
			c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorUnknown, "系統錯誤", nil).Get())
		}
		return
	}

	groupID, code, err := group.Create(&req, owner)
	if err != nil {
		c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorUnknown, "建立群組失敗", nil).Get())
		return
	}

	c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorNone, "OK", map[string]any{
		"group_id": groupID,
		"code":     code,
	}).Get())
}
