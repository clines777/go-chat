package model

type UploadFile struct {
	ID         int    `db:"id" json:"id"`
	Type       int16  `db:"type" json:"type"`
	Path       string `db:"path" json:"path"`
	MD5        string `db:"md5" json:"md5"`
	CreateTime int    `db:"create_time" json:"create_time"`
}
