package funneling

import (
	"context"
	"fmt"
	"github.com/itsnoproblem/prmry/internal/funnel"

	"github.com/itsnoproblem/prmry/internal/components"
	funnelcmp "github.com/itsnoproblem/prmry/internal/components/funnel"
)

//func formatListFunnelsResponse(ctx context.Context, response interface{}) (components.Component, error) {
//	funnels, ok := response.([]funnel.Summary)
//	if !ok {
//		return &components.BaseComponent{}, fmt.Errorf("formatFlowSummaries: failed to cast response")
//	}
//
//	summaries := make([]flowcmp.FlowSummary, 0)
//	for _, flow := range funnels {
//		label := "rule"
//		if len(funnels) > 1 {
//			label = "rules"
//		}
//
//		summaries = append(summaries, funnelcmp.FunnelListView{
//			ID:          flow.ID,
//			Name:        flow.Name,
//			RuleCount:   fmt.Sprintf("%d %s", len(flow.Triggers), label),
//			LastChanged: flow.UpdatedAt.Format("Jan 02, 2006 15:04"),
//		})
//	}
//
//	cmp := flowcmp.FlowsListView{
//		Flows: summaries,
//	}
//	cmp.SetUser(auth.UserFromContext(ctx))
//	cmp.SetTemplates(flowcmp.FlowsListPage(cmp), flowcmp.FlowsList(cmp))
//
//	return &cmp, nil
//}

func formatFunnelBuilderResponse(ctx context.Context, response interface{}) (components.Component, error) {
	res, ok := response.(funnel.Funnel)
	if !ok {
		return &components.BaseComponent{}, fmt.Errorf("formatCreateFunnelResponse: failed to cast response")
	}

	cmp := funnelcmp.FunnelFormView{
		ID:   res.ID,
		Name: res.Name,
		Path: res.Path,
	}
	cmp.SetTemplates(funnelcmp.FunnelBuilderPage(cmp), funnelcmp.FunnelBuilder(cmp))

	return &cmp, nil
}
