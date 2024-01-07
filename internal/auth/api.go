package auth

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

type APIKey struct {
	Name      string
	Key       string
	CreatedAt time.Time
}

func GenerateAPIKey() (string, error) {
	bytes := make([]byte, 16) // 16 bytes = 32 characters in hex
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
