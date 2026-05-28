package model

type Avatar struct {
	ID       int    `db:"id" json:"id"`
	Filename string `db:"filename" json:"filename"`
}
