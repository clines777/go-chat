package log

import (
	"encoding/json"
	"gochat/internal/infra/db"
	"gochat/internal/ws"
	"strconv"
)

const (
	TypeConnErr = iota
	TypeInvalidOpt
)

func SaveLog(logType int32, connID string) error {
	var info json.RawMessage
	if connID != "" {
		sess := ws.GetSocketSession(connID)
		if sess != nil {
			var err error
			userId := strconv.FormatInt(sess.UserID, 10)
			info, err = json.Marshal(
				map[string]string{
					"site_bid": sess.SiteBid,
					"conn_id":  connID,
					"user_id":  userId,
				},
			)
			if err != nil {
				return err
			}
		}
	}

	d, err := db.GetDBConn()
	if err != nil {
		return err
	}

	_, err = d.Builder.
		Insert("log").
		Columns("type", "info").
		Values(logType, info).
		RunWith(d.DB).
		Exec()

	return err
}
