package model

type Config struct {
	ID        int    `db:"id" json:"id"`
	ConfigKey string `db:"config_key" json:"config_key"`
	Value     string `db:"value" json:"value"`
	Remark    string `db:"remark" json:"remark"`
	IsOpen    bool   `db:"is_open" json:"is_open"`
}
