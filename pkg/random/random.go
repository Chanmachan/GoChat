package random

import (
	"crypto/rand"
	"encoding/base64"
)

// ランダムに文字列を生成する
func GenerateRandomString() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	// base64 Encoding
	return base64.URLEncoding.EncodeToString(b)
}
