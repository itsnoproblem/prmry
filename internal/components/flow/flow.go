package flow

import (
	"github.com/itsnoproblem/prmry/internal/components"
	"github.com/itsnoproblem/prmry/internal/flow"
	"sort"
)

type Detail struct {
	ID                  string
	Name                string
	Rules               []RuleView
	RequireAll          bool
	Prompt              string
	PromptArgs          []PromptArg
	SupportedFields     SortedMap
	SupportedConditions SortedMap
	AvailableFlowsByID  SortedMap
	components.BaseComponent
}

func (d *Detail) SetAvalableFlows(flows []flow.Flow) {
	d.AvailableFlowsByID = make(SortedMap, len(flows))
	for _, f := range flows {
		d.AvailableFlowsByID[f.ID] = f.Name
	}
}

type PromptArg struct {
	Source flow.SourceType
	Value  string
}

func (d *Detail) ToFlow() flow.Flow {
	rules := make([]flow.Rule, 0)
	for _, cnd := range d.Rules {
		rules = append(rules, flow.Rule{
			Field: flow.Field{
				Source: flow.SourceType(cnd.Field.Source),
				Value:  cnd.Field.Value,
			},
			Condition: flow.ConditionType(cnd.Condition),
			Value:     cnd.Value,
		})
	}

	promptArgs := make([]flow.Field, 0)
	for _, arg := range d.PromptArgs {
		promptArgs = append(promptArgs, flow.Field{
			Source: arg.Source,
			Value:  arg.Value,
		})
	}

	return flow.Flow{
		ID:         d.ID,
		Name:       d.Name,
		Rules:      rules,
		RequireAll: d.RequireAll,
		Prompt:     d.Prompt,
		PromptArgs: promptArgs,
	}
}

type Field struct {
	Source string
	Value  string
}

type RuleView struct {
	Field
	Condition string
	Value     string
}

type FlowsListView struct {
	Flows []FlowSummary
	components.BaseComponent
}

type FlowSummary struct {
	ID          string
	Name        string
	RuleCount   string
	LastChanged string
}

func NewDetail(flw flow.Flow) Detail {
	rules := make([]RuleView, len(flw.Rules))
	for i, r := range flw.Rules {
		rules[i] = RuleView{
			Field: Field{
				Source: r.Field.Source.String(),
				Value:  r.Field.Value,
			},
			Condition: r.Condition.String(),
			Value:     r.Value,
		}
	}

	promptArgs := make([]PromptArg, 0)
	for _, arg := range flw.PromptArgs {
		promptArgs = append(promptArgs, PromptArg{
			Source: arg.Source,
			Value:  arg.Value,
		})
	}

	return Detail{
		ID:                  flw.ID,
		Name:                flw.Name,
		Rules:               rules,
		RequireAll:          flw.RequireAll,
		Prompt:              flw.Prompt,
		PromptArgs:          promptArgs,
		SupportedFields:     SortedMap(flow.SupportedFields()),
		SupportedConditions: SortedMap(flow.SupportedConditions()),
		AvailableFlowsByID:  nil,
	}
}

type SortedMap map[string]string

func (s SortedMap) Keys() []string {
	keys := make([]string, 0)
	for k, _ := range s {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
