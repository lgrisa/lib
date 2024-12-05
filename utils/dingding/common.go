package dingding

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func NewDingRobotWithSecret(token string, secret string) (*DingRobot, error) {
	if token == "" || secret == "" {
		return nil, fmt.Errorf("no token or no secret")
	}
	return &DingRobot{
		token,
		secret,
	}, nil
}

func sign(timestamp int64, secret string) string {
	data := fmt.Sprintf("%d\n%s", timestamp, secret)

	// HMAC-SHA256 sign
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	signature := h.Sum(nil)

	// Base64 encoding
	encoded := base64.StdEncoding.EncodeToString(signature)
	return encoded
}
