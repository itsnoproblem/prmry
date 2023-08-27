package interaction

import (
	"time"

	openai "github.com/sashabaranov/go-openai"
)

type Moderation struct {
	ID            string
	Model         string
	InteractionID string
	Results       []openai.Result
	CreatedAt     time.Time
}
