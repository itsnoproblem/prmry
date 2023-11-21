package chat

import (
	"strconv"

	"github.com/itsnoproblem/prmry/internal/components"
	"github.com/itsnoproblem/prmry/internal/flow"
	"github.com/itsnoproblem/prmry/internal/interaction"
)

type DetailView struct {
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
	Params       components.SortedMap
}

type Flow struct {
	ID   string
	Name string
}

func NewFlowSelector(flows []flow.Flow, selectedFlow string) FlowSelector {
	flws := make([]Flow, 0)
	params := make(map[string]string)
	for _, flw := range flows {
		flws = append(flws, Flow{
			ID:   flw.ID,
			Name: flw.Name,
		})

		if flw.ID == selectedFlow {
			for _, param := range flw.InputParams {
				params[param.Key] = strconv.FormatBool(param.IsRequired)
			}
		}
	}

	return FlowSelector{
		Flows:        flws,
		SelectedFlow: selectedFlow,
		Params:       components.SortedMap(params),
	}
}

type ControlsView struct {
	FlowSelector FlowSelector
	SelectedFlow flow.Flow
	components.BaseComponent
}

type ResponseView struct {
	Interaction DetailView
	Controls    ControlsView
	components.BaseComponent
}

func NewDetailView(ixn interaction.Interaction) DetailView {
	prompt := ixn.Prompt
	if prompt == "" && len(ixn.Request.Messages) > 0 {
		prompt = ixn.Request.Messages[0].Content
	}

	return DetailView{
		Prompt:       prompt,
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
