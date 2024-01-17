package flow

import (
	"fmt"
	"github.com/itsnoproblem/prmry/internal/components"
	"github.com/itsnoproblem/prmry/internal/flow"
	"strconv"
)

const (
	TabNamePrompt  = "prompt"
	TabNameTrigger = "trigger"
	TabNameInputs  = "input"
	TabNamePreview = "preview"
)

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

type PromptArg struct {
	Source flow.SourceType
	Value  string
}

type FlowSummary struct {
	ID          string
	Name        string
	RuleCount   string
	LastChanged string
}

type InputParam struct {
	Type     flow.ParamType
	Key      string
	Required bool
	Value    string
}

type InputParams []InputParam

func (p InputParams) Map() components.SortedMap {
	m := make(map[string]string)
	for _, param := range p {
		m[param.Key] = param.Key
	}

	return components.SortedMap(m)
}

type Detail struct {
	components.BaseComponent
	ID                  string
	Name                string
	Model               string
	Temperature         string
	Rules               []RuleView
	SelectedTab         string
	RequireAll          bool
	Prompt              string
	PromptArgs          []PromptArg
	SupportedModels     components.SortedMap
	SupportedFields     components.SortedMap
	SupportedConditions components.SortedMap
	AvailableFlowsByID  components.SortedMap
	InputParams         InputParams
}

func (d *Detail) AvailableTags() components.SortedMap {
	tags := make(map[string]string, 0)
	for _, arg := range d.PromptArgs {
		if arg.Source == flow.FieldSourceInputArg {
			tags[arg.Value] = arg.Value
		}
	}
	return tags
}

func (d *Detail) SetAvalableFlows(flows []flow.Flow) {
	d.AvailableFlowsByID = make(components.SortedMap, len(flows))
	for _, f := range flows {
		d.AvailableFlowsByID[f.ID] = f.Name
	}
}

func (d *Detail) ToFlow() flow.Flow {
	rules := make([]flow.Trigger, 0)
	for _, cnd := range d.Rules {
		rules = append(rules, flow.Trigger{
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

	inputParams := make([]flow.InputParam, 0)
	for _, param := range d.InputParams {
		inputParams = append(inputParams, flow.InputParam{
			Type:       param.Type,
			Key:        param.Key,
			IsRequired: param.Required,
		})
	}

	temp, err := strconv.ParseFloat(d.Temperature, 64)
	if err != nil {
		temp = 0
	}

	return flow.Flow{
		ID:          d.ID,
		Name:        d.Name,
		Model:       d.Model,
		Temperature: temp,
		Triggers:    rules,
		RequireAll:  d.RequireAll,
		Prompt:      d.Prompt,
		PromptArgs:  promptArgs,
		InputParams: inputParams,
	}
}

func NewDetail(flw flow.Flow) Detail {
	rules := make([]RuleView, len(flw.Triggers))
	for i, r := range flw.Triggers {
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

	temp := flow.DefaultTemperature
	if flw.Temperature != 0 {
		temp = flw.Temperature
	}

	model := flow.DefaultModel
	if flw.Model != "" {
		model = flw.Model
	}

	inputParams := make([]InputParam, 0)
	for _, param := range flw.InputParams {
		inputParams = append(inputParams, InputParam{
			Type:     param.Type,
			Key:      param.Key,
			Required: param.IsRequired,
		})
	}

	return Detail{
		ID:                  flw.ID,
		Name:                flw.Name,
		Model:               model,
		Temperature:         fmt.Sprintf("%.1f", temp),
		Rules:               rules,
		RequireAll:          flw.RequireAll,
		Prompt:              flw.Prompt,
		PromptArgs:          promptArgs,
		SupportedFields:     components.SortedMap(flow.SupportedFields()),
		SupportedConditions: components.SortedMap(flow.SupportedConditions()),
		SupportedModels:     components.SortedMap(flow.SupportedModels()),
		InputParams:         inputParams,
		AvailableFlowsByID:  nil,
	}
}
