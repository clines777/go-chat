package model

import "time"

type SysLastRead struct {
	Id         int       `db:"column:id" json:"id"`
	UserId     int       `db:"column:user_id" json:"user_id"`
	LastReadId int       `db:"column:last_read_id" json:"last_read_id"`
	UpdateTime time.Time `db:"column:update_time" json:"update_time"`
}
