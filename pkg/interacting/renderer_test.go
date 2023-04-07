package interacting_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	approvals "github.com/approvals/go-approval-tests"
	gogpt "github.com/sashabaranov/go-gpt3"

	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/auth"
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/htmx"
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/interacting"
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/interaction"
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/templates"
)

const (
	testDataDir = "testdata"
)

type RenderFunc func(w http.ResponseWriter, r *http.Request, fullPageTemplate string, fragmentTemplate string, cmp htmx.Component) error

func TestRenderChat(t *testing.T) {
	//r := approvals.UseReporter(reporters.NewGoLandReporter())
	//defer r.Close()
	approvals.UseFolder(testDataDir)

	recorder := httptest.NewRecorder()

	usr := authUser()
	ctx := context.WithValue(context.Background(), auth.ContextKey, &usr)
	reader := bytes.NewReader([]byte(``))
	reqWithUser, err := http.NewRequestWithContext(ctx, http.MethodGet, "/fake", reader)

	if err != nil {
		t.Fatalf("failed to create http request: %s", err)
	}

	tpl, err := templates.Parse()
	if err != nil {
		t.Fatalf("Failed to parse templates: %s", err)
	}
	renderer := htmx.NewRenderer(tpl)
	testInteraction := newInteraction()
	detailView := interacting.DetailView(testInteraction)

	testInteractions := newSummaries(5)
	listView := interacting.ListView(testInteractions)
	listView.SetUser(usr)

	tt := []struct {
		Description      string
		PageTemplate     string
		FragmentTemplate string
		Component        htmx.Component
	}{
		{
			Description:      "HTML fragment of the chat console",
			PageTemplate:     "chat-response.gohtml",
			FragmentTemplate: "chat-response.gohtml",
			Component: &interacting.ChatResponse{
				Interaction: detailView,
			},
		},
		{
			Description:      "Single Interaction",
			PageTemplate:     "page-interaction.gohtml",
			FragmentTemplate: "interaction-single.gohtml",
			Component:        &detailView,
		},
		{
			Description:      "Interaction List",
			PageTemplate:     "page-interactions.gohtml",
			FragmentTemplate: "interactions-list.gohtml",
			Component:        &listView,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Description, func(t *testing.T) {
			err = renderer.RenderComponent(recorder, reqWithUser, tc.PageTemplate, tc.FragmentTemplate, tc.Component)
			if err != nil {
				t.Errorf("ERROR in %s: %s", tc.Description, err)
			}

			approvals.VerifyString(t, recorder.Body.String())
		})
	}
}

func authUser() *auth.User {
	return &auth.User{
		ID:         "abcdefg",
		Name:       "Joey Calzone",
		Nickname:   "jcalzone",
		Email:      "jcal@zone.org",
		AvatarURL:  "http://fake.com/avatar.jpg",
		Provider:   "fakeProv",
		ProviderID: "1234",
	}
}

func newSummaries(num int) []interaction.Summary {
	ixns := make([]interaction.Summary, num)
	for i := 0; i < num; i++ {
		ixns[i] = interaction.Summary{
			ID:             "1234abde",
			Type:           "textresponse",
			Model:          "dabubu",
			Prompt:         "What is this?",
			TokensUsed:     20,
			ResponseLength: 429,
			Error:          "",
			CreatedAt:      time.Date(2022, 12, 12, 0, 0, 0, 0, time.UTC),
		}
	}
	return ixns
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
