package api

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

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

	coverURL := ""
	if g.CoverFileName != "" {
		coverURL = "/static/group-covers/" + g.CoverFileName
	}

	bannedIDs, err := group.GetBannedUserIDs(g.ID)
	if err != nil {
		log.Printf("[GetGroupInfo] GetBannedUserIDs group=%d error: %v", g.ID, err)
		c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorUnknown, "系統錯誤", nil).H())
		return
	}

	c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorNone, "OK", &protocol.GroupInfoResp{
		Title:         g.Title,
		UserTotal:     count,
		Bulletin:      g.Bulletin,
		OwnerUsername: g.OwnerUserName,
		OwnerUserID:   g.OwnerUserID,
		Code:          g.Code,
		Remark:        g.Remark,
		CoverURL:      coverURL,
		PinChatID:     g.PinChatID,
		BannedUserIDs: bannedIDs,
	}).H())
}

// CreateGroup - 創建群組
func CreateGroup(c *gin.Context) {
	var req protocol.CreateGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrWrongParam, "參數錯誤", nil).H())
		return
	}

	userID := c.MustGet("user_id").(int)

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

// UploadGroupCover - 群主上傳群封面
func UploadGroupCover(c *gin.Context) {
	const maxSize = 1024 * 1024 * 2 // 2 MB
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize+1024)

	if err := c.Request.ParseMultipartForm(maxSize); err != nil {
		log.Printf("[UploadGroupCover] upload group cover err: %v", err)
		c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrUploadFileSizeTooLarge, "檔案過大或格式錯誤（上限 2MB）", nil).H())
		return
	}

	groupID, err := strconv.Atoi(c.PostForm("group_id"))
	if err != nil || groupID == 0 {
		c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrWrongParam, "參數錯誤", nil).H())
		return
	}

	userID := c.MustGet("user_id").(int)

	membership, err := group.GetMembership(userID, groupID)
	if err != nil {
		c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrUnauthorized, "無權限", nil).H())
		return
	}
	if membership.RoleType != group.RoleOwner {
		c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrUnauthorized, "只有群主可以更換封面", nil).H())
		return
	}

	file, _, err := c.Request.FormFile("cover")
	if err != nil {
		c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrWrongParam, "請選擇圖片", nil).H())
		return
	}
	defer file.Close()

	// 讀前 512 bytes 偵測 MIME，再接續寫入，避免截斷檔案
	sniff := make([]byte, 512)
	n, _ := file.Read(sniff)
	mimeType := http.DetectContentType(sniff[:n])

	var ext string
	switch mimeType {
	case "image/jpeg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	default:
		c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrWrongParam, "僅支援 JPEG / PNG 格式", nil).H())
		return
	}

	g, err := group.FindByID(groupID)
	if err != nil {
		log.Printf("[UploadGroupCover] FindByID group=%d error: %v", groupID, err)
		c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorUnknown, "系統錯誤", nil).H())
		return
	}

	const dir = "./static/group-covers"
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("[UploadGroupCover] MkdirAll: %v", err)
		c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorUnknown, "系統錯誤", nil).H())
		return
	}

	filename := randHex(16) + ext
	dst, err := os.Create(filepath.Join(dir, filename))
	if err != nil {
		log.Printf("[UploadGroupCover] Create file: %v", err)
		c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorUnknown, "儲存失敗", nil).H())
		return
	}

	if _, err := dst.Write(sniff[:n]); err != nil {
		dst.Close()
		os.Remove(filepath.Join(dir, filename))
		c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorUnknown, "儲存失敗", nil).H())
		return
	}
	if _, err := io.Copy(dst, file); err != nil {
		dst.Close()
		os.Remove(filepath.Join(dir, filename))
		c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorUnknown, "儲存失敗", nil).H())
		return
	}
	dst.Close()

	if err := group.UpdateCoverFilename(groupID, filename); err != nil {
		os.Remove(filepath.Join(dir, filename))
		log.Printf("[UploadGroupCover] UpdateCoverFilename: %v", err)
		c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorUnknown, "更新失敗", nil).H())
		return
	}

	if g.CoverFileName != "" {
		_ = os.Remove(filepath.Join(dir, g.CoverFileName))
	}

	c.JSON(http.StatusOK, protocol.NewApiResponse(protocol.ErrorNone, "OK", map[string]any{
		"url": "/static/group-covers/" + filename,
	}).H())
}

func randHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
