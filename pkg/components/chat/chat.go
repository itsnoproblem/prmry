package chat

import (
	"strconv"

	"github.com/itsnoproblem/prmry/pkg/components"
	"github.com/itsnoproblem/prmry/pkg/interaction"
)

const PromptMaxCharacters = 210

type ChatSummaryView struct {
	ID         string
	Prompt     string
	Type       string
	Date       string
	Model      string
	TokensUsed int
}

type ChatDetailView struct {
	Prompt       string
	PromptHTML   string
	Date         string
	ResponseText string
	ResponseHTML string
	Model        string
	Usage        ChatUsageView
	components.BaseComponent
}

type ChatUsageView struct {
	PromptTokens     string
	CompletionTokens string
	TotalTokens      string
}

type PersonaSelector struct {
	ID   string
	Name string
}

type ChatControlsView struct {
	Personas []PersonaSelector
	components.BaseComponent
}

type ChatResponseView struct {
	Interaction ChatDetailView
	Controls    ChatControlsView
	components.BaseComponent
}

type InteractionListView struct {
	Interactions []ChatSummaryView
	components.BaseComponent
}

func NewInteractionListView(summaries []interaction.Summary) InteractionListView {
	interactions := make([]ChatSummaryView, len(summaries))
	for i, s := range summaries {
		interactions[i] = ChatSummaryView{
			ID:         s.ID,
			Prompt:     components.TrimWordsToMaxCharacters(PromptMaxCharacters, s.Prompt),
			Type:       s.Type,
			Date:       s.CreatedAt.Format("Jan 2, 2006 3:04pm"),
			Model:      s.Model,
			TokensUsed: s.TokensUsed,
		}
	}

	return InteractionListView{
		Interactions: interactions,
	}
}

func NewChatDetailView(ixn interaction.Interaction) ChatDetailView {
	return ChatDetailView{
		Prompt:       ixn.Request.Prompt,
		PromptHTML:   ixn.PromptHTML(),
		Date:         ixn.CreatedAt.Format("Monday, January 2, 2006 at 3:04pm"),
		ResponseText: ixn.ResponseText(),
		ResponseHTML: ixn.ResponseHTML(),
		Model:        ixn.Response.Model,
		Usage: ChatUsageView{
			PromptTokens:     strconv.Itoa(ixn.Response.Usage.PromptTokens),
			CompletionTokens: strconv.Itoa(ixn.Response.Usage.CompletionTokens),
			TotalTokens:      strconv.Itoa(ixn.Response.Usage.TotalTokens),
		},
	}
}
