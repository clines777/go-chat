package group

import (
	"gochat/internal/infra/db"
	"gochat/internal/protocol"
)

const getGroupsOfUserSQL = `
SELECT cg.id, cg.title, cg.code, cg.open_join, cg.join_user_level,
       COALESCE(lr.content, '') AS last_msg,
       COALESCE(lr.create_time, '1970-01-01 00:00:00+00') AS last_msg_time
FROM chat_group cg
JOIN group_user gu ON gu.group_id = cg.id
    AND gu.user_id = $1
    AND gu.deleted = false
LEFT JOIN LATERAL (
    SELECT content, create_time
    FROM chat_record
    WHERE group_id = cg.id AND deleted = false
    ORDER BY id DESC LIMIT 1
) lr ON true
WHERE cg.site_bid = $2
  AND cg.is_dismiss = false`

func GetGroupsOfUser(userID int64, siteBid string) ([]protocol.DisplayUserGroup, error) {
	d, err := db.GetDBConn()
	if err != nil {
		return nil, err
	}

	rows := make([]protocol.DisplayUserGroup, 0)
	if err := d.DB.Select(&rows, getGroupsOfUserSQL, userID, siteBid); err != nil {
		return nil, err
	}

	return rows, nil
}
