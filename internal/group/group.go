package group

import (
	"gochat/internal/infra/db"
	"gochat/internal/protocol"
)

func GetGroupsOfUser(userID int64, siteBid string) ([]protocol.DisplayUserGroup, error) {
	d, err := db.GetDBConn()
	if err != nil {
		return nil, err
	}

	querySql, args, _ := d.Builder.
		Select("g.title", "g.code", "g.open_join", "g.join_user_level", "g.id").
		From("chat_group g").
		InnerJoin("group_user gu ON g.id = gu.group_id").
		Where("gu.user_id = ?", userID).
		Where("g.site_bid = ?", siteBid).
		ToSql()

	rows := make([]protocol.DisplayUserGroup, 0)
	if err := d.DB.Select(&rows, querySql, args...); err != nil {
		return rows, nil
	}

	return rows, nil
}
