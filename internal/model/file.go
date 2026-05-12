package model

type UploadFile struct {
	ID         int64  `db:"id" json:"id"`
	Type       int16  `db:"type" json:"type"`
	Path       string `db:"path" json:"path"`
	SiteBid    string `db:"site_bid" json:"site_bid"`
	MD5        string `db:"md5" json:"md5"`
	CreateTime int64  `db:"create_time" json:"create_time"`
}
