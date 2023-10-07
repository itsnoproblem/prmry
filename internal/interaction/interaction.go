package interaction

import (
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

type Interaction struct {
	ID               string
	FlowID           string
	FlowName         string
	Request          openai.ChatCompletionRequest
	Response         openai.ChatCompletionResponse
	Type             string
	Model            string
	Prompt           string
	Completion       string
	TokensPrompt     int
	TokensCompletion int
	CreatedAt        time.Time
	UserID           string
}

type Summary struct {
	ID             string
	FlowID         string
	FlowName       string
	Type           string
	Model          string
	Prompt         string
	TokensUsed     int
	ResponseLength int
	CreatedAt      time.Time
	UserID         string
}

func (s Summary) PromptFragment(charLimit int) string {
	if len(s.Prompt) <= charLimit {
		return s.Prompt
	}

	for i, char := range s.Prompt {
		if i > charLimit && char == ' ' {
			return s.Prompt[:i] + "..."
		}
	}

	return s.Prompt
}

// PromptHTML returns the prompt as HTML
func (ixn Interaction) PromptHTML() string {
	if ixn.Prompt == "" && len(ixn.Request.Messages) < 1 {
		return ""
	}

	prompt := ixn.Prompt
	if prompt == "" {
		prompt = ixn.Request.Messages[0].Content
	}
	return textToHTML(prompt)
}

// ResponseText returns the text of the final element of the Prompt.Choices slice
func (ixn Interaction) ResponseText() string {
	text := ixn.Completion
	if text == "" && len(ixn.Response.Choices) > 0 {
		for _, res := range ixn.Response.Choices {
			text = res.Message.Content
		}
	}
	return text
}

// ResponseHTML returns the final element of the Prompt.Choices slice as HTML
func (ixn Interaction) ResponseHTML() string {
	text := textToHTML(ixn.ResponseText())
	return text
}

func textToHTML(text string) string {
	htm := strings.Trim(text, "\r\n")
	htm = strings.Trim(text, "\n")
	htm = strings.ReplaceAll(htm, "\r\n", "<br/>")
	htm = strings.ReplaceAll(htm, "\n", "<br/>")
	return htm
}
