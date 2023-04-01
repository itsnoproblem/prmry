package interacting

import (
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/htmx"
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/interaction"
	"html/template"
	"strings"
)

const PromptMaxCharacters = 140

type ListViewModel struct {
	Interactions []SummaryViewModel
	htmx.BaseComponent
}

type SummaryViewModel struct {
	ID         string
	Prompt     string
	Type       string
	Date       string
	Model      string
	TokensUsed int
}

type DetailViewModel struct {
	Prompt       string
	PromptHTML   template.HTML
	Date         string
	ResponseText string
	ResponseHTML template.HTML
	Model        string
	Usage        UsageViewModel
	htmx.BaseComponent
}

type UsageViewModel struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

type ChatResponse struct {
	Interaction DetailViewModel
	htmx.BaseComponent
}

type InteractionList struct {
	Interactions []interaction.Summary
	htmx.BaseComponent
}

func ListView(summaries []interaction.Summary) ListViewModel {
	interactions := make([]SummaryViewModel, len(summaries))
	for i, s := range summaries {
		interactions[i] = SummaryViewModel{
			ID:         s.ID,
			Prompt:     trimWordsToMaxCharacters(PromptMaxCharacters, s.Prompt),
			Type:       s.Type,
			Date:       s.CreatedAt.Format("Jan 2, 2006 3:04pm"),
			Model:      s.Model,
			TokensUsed: s.TokensUsed,
		}
	}

	return ListViewModel{
		Interactions: interactions,
	}
}

func DetailView(ixn interaction.Interaction) DetailViewModel {
	return DetailViewModel{
		Prompt:       ixn.Request.Prompt,
		PromptHTML:   template.HTML(ixn.PromptHTML()),
		Date:         ixn.CreatedAt.Format("Monday, January 2, 2006 at 3:04pm"),
		ResponseText: ixn.ResponseText(),
		ResponseHTML: template.HTML(ixn.ResponseHTML()),
		Model:        ixn.Response.Model,
		Usage: UsageViewModel{
			PromptTokens:     ixn.Response.Usage.PromptTokens,
			CompletionTokens: ixn.Response.Usage.CompletionTokens,
			TotalTokens:      ixn.Response.Usage.TotalTokens,
		},
	}
}

func trimWordsToMaxCharacters(maxChars int, text string) string {
	var (
		trimmed = ""
		strlen  = 0
	)
	fields := strings.Fields(text)
	for i, f := range fields {
		strlen += len(f)
		if strlen > maxChars {
			trimmed = strings.TrimSuffix(trimmed, " ") + "..."
			break
		}

		if i < len(fields)-1 {
			trimmed += f + " "
		}
	}

	return trimmed
}
