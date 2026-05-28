package model

import (
	"encoding/json"
)

type Log struct {
	ID   int             `db:"id,omitempty" json:"id"`
	Type int             `db:"type,omitempty" json:"type"`
	Info json.RawMessage `db:"info,omitempty" json:"info,omitempty"`
}
