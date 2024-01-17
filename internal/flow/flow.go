package flow

import (
	"github.com/sashabaranov/go-openai"
	"strings"
	"time"
)

const (
	ParamTypeString ParamType = "string"
	ParamTypeNumber ParamType = "number"

	DefaultTemperature = 0.5
	DefaultModel       = openai.GPT3Dot5Turbo
	DefaultMaxTokens   = 4000
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
	Triggers    []Trigger
	RequireAll  bool
	Prompt      string
	PromptArgs  []Field
	CreatedAt   time.Time
	UpdatedAt   time.Time
	InputParams []InputParam
	Model       string
	Temperature float64
}

func (f Flow) TriggerDescription() string {
	if len(f.Triggers) == 0 {
		return "Executes on any input"
	}

	description := "Executes when "
	var ruleDescriptions []string
	for _, rule := range f.Triggers {
		ruleDescriptions = append(ruleDescriptions, rule.String())
	}

	operator := "OR"
	if f.RequireAll {
		operator = "AND"
	}

	description += strings.Join(ruleDescriptions, "\n"+operator+" ")
	return description
}

func SupportedModels() map[string]string {
	return map[string]string{
		openai.GPT3Dot5Turbo: "GPT 3.5-turbo",
		openai.GPT4:          "GPT 4",
	}
}
