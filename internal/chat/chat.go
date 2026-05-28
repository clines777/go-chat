package chat

import (
	"encoding/json"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"gochat/internal/infra/db"
	infranats "gochat/internal/infra/nats"
	"gochat/internal/model"
	"gochat/internal/protocol"
)

const TypeText = int16(1)

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

func Publish(event *protocol.CastChatEvent) error {
	nc, err := infranats.GetNats()
	if err != nil {
		return err
	}
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return nc.Publish(fmt.Sprintf("chat.group.%d", event.GroupId), data)
}
