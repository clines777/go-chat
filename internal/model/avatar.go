package model

type Avatar struct {
	ID       int64  `db:"id" json:"id"`
	Filename string `db:"filename" json:"filename"`
}
