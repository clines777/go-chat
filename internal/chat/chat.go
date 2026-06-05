package chat

import (
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"gochat/internal/infra/db"
	"gochat/internal/model"
	"gochat/internal/protocol"
	"time"
)

const TypeText = 1

// SaveRecord 寫入群組用戶的發送的聊天訊息.
func SaveRecord(groupID int, userID int, content string) (*model.ChatRecord, error) {
	d, err := db.GetDBConn()
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	var id int

	err = d.Builder.
		Insert("chat_record").
		Columns("group_id", "user_id", "type", "content", "create_time", "update_time").
		Values(groupID, userID, TypeText, content, now, now).
		Suffix("RETURNING id").
		RunWith(d.DB).
		QueryRow().
		Scan(&id)
	if err != nil {
		return nil, err
	}

	return &model.ChatRecord{
		ID:         id,
		GroupID:    groupID,
		UserID:     userID,
		Type:       TypeText,
		Content:    content,
		CreateTime: int(now),
	}, nil
}

// ExistsInGroup 檢查一筆未刪除的聊天訊息是否存在於指定群組。
func ExistsInGroup(chatID, groupID int) (bool, error) {
	d, err := db.GetDBConn()
	if err != nil {
		return false, err
	}
	querySql, args, _ := d.Builder.
		Select("1").
		From("chat_record").
		Where(sq.Eq{"id": chatID, "group_id": groupID, "deleted": false}).
		Limit(1).ToSql()
	var one int
	err = d.DB.Get(&one, querySql, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// MarkDeleted 軟刪除一筆發言, 以 group_id 限定範圍確保訊息屬於該群。回傳受影響筆數。
func MarkDeleted(chatID, groupID int) (int, error) {
	d, err := db.GetDBConn()
	if err != nil {
		return 0, err
	}
	r, err := d.Builder.
		Update("chat_record").
		Set("deleted", true).
		Set("update_time", time.Now().Unix()).
		Where(sq.Eq{"id": chatID, "group_id": groupID, "deleted": false}).
		RunWith(d.DB).Exec()
	if err != nil {
		return 0, err
	}
	n, err := r.RowsAffected()
	return int(n), err
}

// FindInGroup 取得單筆未刪除發言 (含作者 username/avatar), 供置頂訊息展示。找不到回傳 sql.ErrNoRows。
func FindInGroup(chatID, groupID int) (*protocol.ChatInfo, error) {
	d, err := db.GetDBConn()
	if err != nil {
		return nil, err
	}

	querySql, args, _ := d.Builder.
		Select(
			"cr.id", "cr.user_id", "u.username AS username",
			"COALESCE('/static/avatars/' || av.filename, '') AS avatar_url",
			"cr.content", "cr.create_time",
		).
		From("chat_record cr").
		Join(`"user" u ON u.id = cr.user_id`).
		LeftJoin("avatar av ON av.id = u.avatar_id").
		Where(sq.Eq{"cr.id": chatID, "cr.group_id": groupID, "cr.deleted": false}).
		Limit(1).ToSql()

	var row protocol.ChatInfo
	if err := d.DB.Get(&row, querySql, args...); err != nil {
		return nil, err
	}
	return &row, nil
}

// GetRecentChats 取得指定群組最新 limit 筆聊天訊息 (進群初次載入)。
// hasMore 表示這批之前是否還有更舊訊息可供向上分頁載入。
func GetRecentChats(groupID int, limit int) ([]protocol.ChatInfo, bool, error) {
	return GetHistoryChats(groupID, 0, limit)
}

// GetHistoryChats 取得 beforeID 之前 (更舊) 的最多 limit 筆訊息, 依 id 正序回傳。
// beforeID = 0 表示從最新一筆往前取; 走 (group_id, id DESC) WHERE deleted=false 的 partial index。
// hasMore 表示在這批之前是否還有更舊訊息。
func GetHistoryChats(groupID, beforeID, limit int) ([]protocol.ChatInfo, bool, error) {
	d, err := db.GetDBConn()
	if err != nil {
		return nil, false, err
	}

	cond := sq.And{sq.Eq{"cr.group_id": groupID, "cr.deleted": false}}
	if beforeID > 0 {
		cond = append(cond, sq.Lt{"cr.id": beforeID})
	}

	// 多取一筆以判斷是否還有更舊訊息。
	sub := d.Builder.
		Select(
			"cr.id", "cr.user_id", "u.username AS username",
			"COALESCE('/static/avatars/' || av.filename, '') AS avatar_url",
			"cr.content", "cr.create_time",
		).
		From("chat_record cr").
		Join(`"user" u ON u.id = cr.user_id`).
		LeftJoin("avatar av ON av.id = u.avatar_id").
		Where(cond).
		OrderBy("cr.id DESC").
		Limit(uint64(limit + 1))

	querySql, args, _ := d.Builder.
		Select("*").
		FromSelect(sub, "t").
		OrderBy("id ASC").
		ToSql()

	var rows []protocol.ChatInfo
	if err := d.DB.Select(&rows, querySql, args...); err != nil {
		return nil, false, err
	}

	chats, hasMore := trimHistoryPage(rows, limit)
	return chats, hasMore, nil
}

// trimHistoryPage 處理「多取一筆」結果。rows 為 id 正序且最多 limit+1 筆:
// 若超過 limit 代表這批之前還有更舊訊息, 砍掉最舊 (開頭) 那筆並回報 hasMore=true。
func trimHistoryPage(rows []protocol.ChatInfo, limit int) ([]protocol.ChatInfo, bool) {
	if len(rows) > limit {
		return rows[len(rows)-limit:], true
	}
	return rows, false
}
