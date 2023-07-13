package flow

import (
	"time"
)

type Flow struct {
	ID         string
	UserID     string
	Name       string
	Rules      []Rule
	RequireAll bool
	Prompt     string
	PromptArgs []Field
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
