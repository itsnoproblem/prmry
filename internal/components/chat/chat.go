package chat

import (
	"github.com/itsnoproblem/prmry/internal/flow"
	"strconv"

	"github.com/itsnoproblem/prmry/internal/components"
	"github.com/itsnoproblem/prmry/internal/interaction"
)

const PromptMaxCharacters = 210

type ChatSummaryView struct {
	ID         string
	Prompt     string
	Type       string
	Date       string
	Model      string
	FlowID     string
	FlowName   string
	TokensUsed int
}

type ChatDetailView struct {
	Prompt       string
	PromptHTML   string
	Date         string
	ResponseText string
	ResponseHTML string
	FlowID       string
	FlowName     string
	Model        string
	Usage        ChatUsageView
	components.BaseComponent
}

type ChatUsageView struct {
	PromptTokens     string
	CompletionTokens string
	TotalTokens      string
}

type FlowSelector struct {
	Flows        []Flow
	SelectedFlow string
}

type Flow struct {
	ID   string
	Name string
}

func NewFlowSelector(flows []flow.Flow, selectedFlow string) FlowSelector {
	flws := make([]Flow, 0)
	for _, flw := range flows {
		flws = append(flws, Flow{
			ID:   flw.ID,
			Name: flw.Name,
		})
	}

	return FlowSelector{
		Flows:        flws,
		SelectedFlow: selectedFlow,
	}
}

type ChatControlsView struct {
	FlowSelector FlowSelector
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
			FlowID:     s.FlowID,
			FlowName:   s.FlowName,
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
		FlowID:       ixn.FlowID,
		FlowName:     ixn.FlowName,
		Model:        ixn.Response.Model,
		Usage: ChatUsageView{
			PromptTokens:     strconv.Itoa(ixn.Response.Usage.PromptTokens),
			CompletionTokens: strconv.Itoa(ixn.Response.Usage.CompletionTokens),
			TotalTokens:      strconv.Itoa(ixn.Response.Usage.TotalTokens),
		},
	}
}
