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

// ExistsInGroup 檢查一筆未刪除的發言是否存在於指定群組。
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

func GetRecentChats(groupID int, limit int) ([]protocol.ChatInfo, error) {
	d, err := db.GetDBConn()
	if err != nil {
		return nil, err
	}

	sub := d.Builder.
		Select(
			"cr.id", "cr.user_id", "u.username AS username",
			"COALESCE('/static/avatars/' || av.filename, '') AS avatar_url",
			"cr.content", "cr.create_time",
		).
		From("chat_record cr").
		Join(`"user" u ON u.id = cr.user_id`).
		LeftJoin("avatar av ON av.id = u.avatar_id").
		Where(sq.Eq{"cr.group_id": groupID, "cr.deleted": false}).
		OrderBy("cr.create_time DESC").
		Limit(uint64(limit))

	querySql, args, _ := d.Builder.
		Select("*").
		FromSelect(sub, "t").
		OrderBy("create_time ASC").
		ToSql()

	var rows []protocol.ChatInfo
	if err := d.DB.Select(&rows, querySql, args...); err != nil {
		return nil, err
	}
	return rows, nil
}
