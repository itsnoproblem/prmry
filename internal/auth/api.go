package auth

import "time"

type APIKey struct {
	Name      string
	Key       string
	CreatedAt time.Time
}
