package interacting_test

import (
	"bytes"
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/interacting"
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/interaction"
	"io"
	"testing"
	"time"

	approvals "github.com/approvals/go-approval-tests"
	"github.com/approvals/go-approval-tests/reporters"
	gogpt "github.com/sashabaranov/go-gpt3"
	"github.com/stretchr/testify/require"
)

const (
	testDataDir = "testdata"
)

func TestRenderChat(t *testing.T) {
	r := approvals.UseReporter(reporters.NewGoLandReporter())
	defer r.Close()
	approvals.UseFolder(testDataDir)

	renderer, err := interacting.NewRenderer()
	require.NoError(t, err)

	tt := []struct {
		Description string
		RenderFn    func(w io.Writer) error
	}{
		{
			"HTML fragment of the chat console",
			renderer.RenderChatConsole,
		},
		{
			"HTML document of the chat console",
			renderer.RenderChatPage,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Description, func(t *testing.T) {
			buf := bytes.Buffer{}
			err = tc.RenderFn(&buf)
			if err != nil {
				t.Errorf("ERROR in %s: %s", tc.Description, err)
			}

			approvals.VerifyString(t, buf.String())
		})
	}
}

func TestRenderInteraction(t *testing.T) {
	r := approvals.UseReporter(reporters.NewGoLandReporter())
	defer r.Close()
	approvals.UseFolder(testDataDir)

	renderer, err := interacting.NewRenderer()
	require.NoError(t, err)

	testInteraction := newInteraction()

	tt := []struct {
		Description string
		Interaction interaction.Interaction
		RenderFn    func(w io.Writer, ixn interaction.Interaction) error
	}{
		{
			Description: "HTML fragment of a single interaction",
			Interaction: testInteraction,
			RenderFn:    renderer.RenderInteraction,
		},
		{
			Description: "HTML document of a single interaction",
			Interaction: testInteraction,
			RenderFn:    renderer.RenderInteractionPage,
		},
		{
			Description: "HTML fragment of a chat interaction",
			Interaction: testInteraction,
			RenderFn:    renderer.RenderChatResponse,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Description, func(t *testing.T) {
			buf := bytes.Buffer{}
			err = tc.RenderFn(&buf, tc.Interaction)
			if err != nil {
				t.Errorf("ERROR in %s: %s", tc.Description, err)
			}

			approvals.VerifyString(t, buf.String())
		})
	}
}

func newInteraction() interaction.Interaction {
	model := gogpt.GPT3TextDavinci003
	return interaction.Interaction{
		ID: "123xyz",
		Request: gogpt.CompletionRequest{
			Model:     model,
			Prompt:    "What is this thing supposed to be?",
			MaxTokens: 4444,
		},
		Response: gogpt.CompletionResponse{
			ID:      "bz321",
			Object:  "",
			Created: uint64(time.Now().UnixNano()),
			Model:   model,
			Choices: []gogpt.CompletionChoice{
				{
					Text: "Well, nobody is quite sure what anything is <i>supposed</i>? to do.",
				},
			},
			Usage: gogpt.Usage{
				PromptTokens:     1001,
				CompletionTokens: 2332,
				TotalTokens:      3333,
			},
		},
		CreatedAt: time.Unix(0, 0),
	}
}
func newSummary() interaction.Summary {
	return interaction.Summary{
		ID:        "xy123",
		Type:      "chat",
		Model:     "davinci-002",
		Prompt:    "What is this thing?",
		CreatedAt: time.Unix(0, 0),
	}
}
