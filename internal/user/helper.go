package user

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func genUserCode(username string) string {
	hash := sha256.Sum256([]byte(strings.ToLower(username)))
	return strings.ToUpper(hex.EncodeToString(hash[:])[:9])
}
