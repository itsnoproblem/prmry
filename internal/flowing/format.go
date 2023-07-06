package flowing

import (
	"context"
	"fmt"
	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/itsnoproblem/prmry/internal/components"
	"github.com/itsnoproblem/prmry/internal/components/success"

	flowcmp "github.com/itsnoproblem/prmry/internal/components/flow"
)

func formatFlowSummaries(ctx context.Context, response interface{}) (components.Component, error) {
	res, ok := response.(listFlowsResponse)
	if !ok {
		return &components.BaseComponent{}, fmt.Errorf("formatFlowSummaries: failed to cast response")
	}

	summaries := make([]flowcmp.FlowSummary, 0)
	for _, flow := range res.Summaries {
		label := "rule"
		if flow.RuleCount > 1 {
			label = "rules"
		}

		summaries = append(summaries, flowcmp.FlowSummary{
			ID:          flow.ID,
			Name:        flow.Name,
			RuleCount:   fmt.Sprintf("%d %s", flow.RuleCount, label),
			LastChanged: flow.LastChanged.Format("Jan 02, 2006 15:04"),
		})
	}

	cmp := flowcmp.FlowsListView{
		Flows: summaries,
	}
	cmp.SetUser(auth.UserFromContext(ctx))
	cmp.SetTemplates(flowcmp.FlowsListPage(cmp), flowcmp.FlowsList(cmp))

	return &cmp, nil
}

func formatFlowBuilderResponse(ctx context.Context, response interface{}) (components.Component, error) {
	resp, ok := response.(flowBuilderResponse)
	if !ok {
		return &components.BaseComponent{}, fmt.Errorf("formatFlowBuilderResponse: failed to parse response")
	}
	resp.Form.SetUser(auth.UserFromContext(ctx))

	fullPage := flowcmp.FlowBuilderPage(resp.Form)
	fragment := flowcmp.FlowBuilder(resp.Form)
	resp.Form.SetTemplates(fullPage, fragment)

	return &resp.Form, nil
}

func formatSuccessMessageResponse(ctx context.Context, response interface{}) (components.Component, error) {
	cmp, ok := response.(success.SuccessView)
	if !ok {
		return &components.BaseComponent{}, fmt.Errorf("formatSuccessMessageResponse: failed to parse response")
	}
	cmp.SetUser(auth.UserFromContext(ctx))

	fragment := success.Success(cmp)
	fullPage := success.SuccessPage(cmp)
	cmp.SetTemplates(fullPage, fragment)

	return &cmp, nil
}

func formatRedirectResponse(ctx context.Context, response interface{}) (components.Component, error) {
	return &components.BaseComponent{}, nil
}
