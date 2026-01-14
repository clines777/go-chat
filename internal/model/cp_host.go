package model

type CpHost struct {
	ID        int64  `db:"id" json:"id"`
	SiteBid   string `db:"site_bid" json:"site_bid"`
	APIHost   string `db:"api_host" json:"api_host"`
	IsEnabled int16  `db:"is_enabled" json:"is_enabled"`
}
