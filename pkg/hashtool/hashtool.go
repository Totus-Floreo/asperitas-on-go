package hashtool

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func CalculateMD5Hash(userID int) string {
	hash := md5.Sum([]byte(fmt.Sprintf("%d", userID)))
	hashString := hex.EncodeToString(hash[:])
	return hashString
}
