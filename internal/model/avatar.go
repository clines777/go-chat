package model

type Avatar struct {
	ID       int64  `db:"id" json:"id"`
	Type     int16  `db:"type" json:"type"`
	Filename string `db:"filename" json:"filename"`
}
