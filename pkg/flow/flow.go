package flow

import (
	"fmt"
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

func (f Flow) Respond(prompt string) (string, error) {
	for _, cond := range f.Rules {
		matches, err := cond.Matches(prompt)
		if err != nil {
			return "", fmt.Errorf("Flow.Respond: %s", err)
		}

		if matches {
			return f.Prompt, nil
		}
	}

	return "", nil
}
