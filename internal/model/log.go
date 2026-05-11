package model

import (
	"encoding/json"
)

type Log struct {
	ID   int64           `db:"id,omitempty" json:"id"`
	Type int32           `db:"type,omitempty" json:"type"`
	Info json.RawMessage `db:"info,omitempty" json:"info,omitempty"`
}
