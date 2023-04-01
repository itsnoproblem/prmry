package interaction

import (
	"strings"
	"time"

	gogpt "github.com/sashabaranov/go-gpt3"
)

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

type Interaction struct {
	ID        string
	Request   gogpt.CompletionRequest
	Response  gogpt.CompletionResponse
	Error     string
	CreatedAt time.Time
}

// PromptHTML returns the prompt as HTML
func (ixn Interaction) PromptHTML() string {
	return textToHTML(ixn.Request.Prompt)
}

// ResponseText returns the text of the final element of the Response.Choices slice
func (ixn Interaction) ResponseText() string {
	var text string
	for _, res := range ixn.Response.Choices {
		text = res.Text
	}
	return text
}

// ResponseHTML returns the final element of the Response.Choices slice as HTML
func (ixn Interaction) ResponseHTML() string {
	return textToHTML(ixn.ResponseText())
}

func textToHTML(text string) string {
	htm := strings.Trim(text, "\r\n")
	htm = strings.Trim(text, "\n")
	htm = strings.ReplaceAll(htm, "\r\n", "<br/>")
	htm = strings.ReplaceAll(htm, "\n", "<br/>")
	return htm
}
