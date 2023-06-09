package flow

import "github.com/itsnoproblem/prmry/pkg/components"

type Detail struct {
	ID                  string
	Name                string
	Conditions          []ConditionView
	RequireAll          bool
	Response            string
	SupportedFields     []string
	SupportedConditions []string
	components.BaseComponent
}

type ConditionView struct {
	Type  string
	Field string
	Value string
}

type FlowsListView struct {
	Flows []FlowSummary
	components.BaseComponent
}

type FlowSummary struct {
	ID   string
	Name string
}
