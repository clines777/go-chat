package user

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

func genUserCode(siteBid string, extMemberID int32, username string) string {
	input := strings.ToLower(siteBid) + fmt.Sprintf("%d", extMemberID) + strings.ToLower(username)
	hash := sha256.Sum256([]byte(input))
	return strings.ToUpper(hex.EncodeToString(hash[:])[:9])
}
