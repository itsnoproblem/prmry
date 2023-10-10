package flow

import (
	"strings"
	"time"
)

const (
	ParamTypeString ParamType = "string"
	ParamTypeNumber ParamType = "number"
)

type ParamType string
type InputParam struct {
	Type       ParamType
	Key        string
	IsRequired bool
}

type Flow struct {
	ID          string
	UserID      string
	Name        string
	Rules       []Rule
	RequireAll  bool
	Prompt      string
	PromptArgs  []Field
	CreatedAt   time.Time
	UpdatedAt   time.Time
	InputParams []InputParam
}

func (f Flow) TriggerDescription() string {
	if len(f.Rules) == 0 {
		return "Executes on any input"
	}

	description := "Executes when "
	var ruleDescriptions []string
	for _, rule := range f.Rules {
		ruleDescriptions = append(ruleDescriptions, rule.String())
	}

	operator := "OR"
	if f.RequireAll {
		operator = "AND"
	}

	description += strings.Join(ruleDescriptions, "\n"+operator+" ")
	return description
}
