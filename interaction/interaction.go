package interaction

import (
	"time"
	
	gogpt "github.com/sashabaranov/go-gpt3"
)

type Interaction struct {
	ID        string
	Request   gogpt.CompletionRequest
	Response  gogpt.CompletionResponse
	Error     string
	CreatedAt time.Time
}

type Summary struct {
	ID             string
	Type           string
	Model          string
	Prompt         string
	TokensUsed     int
	ResponseLength int
	Error          string
	CreatedAt      time.Time
}
