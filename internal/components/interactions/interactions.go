package interactions

import (
	"github.com/itsnoproblem/prmry/internal/components"
	"github.com/itsnoproblem/prmry/internal/interaction"
)

const PromptSummaryMaxCharacters = 210

type ListView struct {
	Interactions []SummaryView
	components.BaseComponent
}

type SummaryView struct {
	ID         string
	Prompt     string
	Type       string
	Date       string
	Model      string
	FlowID     string
	FlowName   string
	TokensUsed int
}

func NewListView(summaries []interaction.Summary) ListView {
	interactions := make([]SummaryView, len(summaries))
	for i, s := range summaries {
		interactions[i] = SummaryView{
			ID:         s.ID,
			Prompt:     components.TrimWordsToMaxCharacters(PromptSummaryMaxCharacters, s.Prompt),
			Type:       s.Type,
			FlowID:     s.FlowID,
			FlowName:   s.FlowName,
			Date:       s.CreatedAt.Format("Jan 2, 2006 3:04pm"),
			Model:      s.Model,
			TokensUsed: s.TokensUsed,
		}
	}

	return ListView{
		Interactions: interactions,
	}
}
