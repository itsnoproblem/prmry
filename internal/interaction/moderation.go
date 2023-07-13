package interaction

import (
	"time"

	gogpt "github.com/sashabaranov/go-gpt3"
)

type Moderation struct {
	ID            string
	Model         string
	InteractionID string
	Results       []gogpt.Result
	CreatedAt     time.Time
}
