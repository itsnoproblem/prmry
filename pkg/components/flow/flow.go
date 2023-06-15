package flow

import (
	"github.com/itsnoproblem/prmry/pkg/components"
	"github.com/itsnoproblem/prmry/pkg/flow"
)

type Detail struct {
	ID                  string
	Name                string
	Rules               []RuleView
	RequireAll          bool
	Prompt              string
	PromptArgs          []string
	SupportedFields     []string
	SupportedConditions []string
	components.BaseComponent
}

func (d Detail) ToFlow() flow.Flow {
	rules := make([]flow.Rule, 0)
	for _, cnd := range d.Rules {
		rules = append(rules, flow.Rule{
			Field:     flow.FieldSource(cnd.Field),
			Condition: flow.ConditionType(cnd.Condition),
			Value:     cnd.Value,
		})
	}

	return flow.Flow{
		ID:         d.ID,
		Name:       d.Name,
		Rules:      rules,
		RequireAll: d.RequireAll,
		Prompt:     d.Prompt,
		PromptArgs: d.PromptArgs,
	}
}

type RuleView struct {
	Field     string
	Condition string
	Value     string
}

type FlowsListView struct {
	Flows []FlowSummary
	components.BaseComponent
}

type FlowSummary struct {
	ID   string
	Name string
}
