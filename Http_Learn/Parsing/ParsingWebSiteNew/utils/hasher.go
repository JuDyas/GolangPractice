package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func HashMD5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
