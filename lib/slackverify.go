package lib

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// SlackVerify verifies a Slack request
func SlackVerify(body []byte, secret string, timestamp string, signature string) bool {
	hash := hmac.New(sha256.New, []byte(secret))
	base := fmt.Sprintf("v0:%s:%s", timestamp, string(body))
	hash.Write([]byte(base))

	return hmac.Equal([]byte(fmt.Sprintf("v0=%s", hex.EncodeToString(hash.Sum(nil)))), []byte(signature))
}
