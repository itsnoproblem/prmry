package funneling

import (
	"context"
	"fmt"
	"github.com/a-h/templ"
	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/itsnoproblem/prmry/internal/funnel"

	"github.com/itsnoproblem/prmry/internal/components"
	funnelcmp "github.com/itsnoproblem/prmry/internal/components/funnel"
)

func formatRedirectResponse(ctx context.Context, response interface{}) (components.Component, error) {
	cmp := components.BaseComponent{}
	cmp.SetUser(auth.UserFromContext(ctx))
	cmp.SetTemplates(templ.NopComponent, templ.NopComponent)

	return &cmp, nil
}

func formatFunnelBuilderResponse(ctx context.Context, response interface{}) (components.Component, error) {
	res, ok := response.(funnel.WithFlows)
	if !ok {
		return &components.BaseComponent{}, fmt.Errorf("formatFunnelBuilderResponse: failed to cast response")
	}

	cmp := funnelcmp.FunnelFormView{
		FunnelView: funnelcmp.FunnelView{
			ID:   res.ID,
			Name: res.Name,
			Path: res.Path,
		},
	}

	for _, f := range res.Flows {
		cmp.Flows = append(cmp.Flows, funnelcmp.FunnelFlowView{
			ID:   f.ID,
			Name: f.Name,
		})
	}

	cmp.SetUser(auth.UserFromContext(ctx))
	cmp.SetTemplates(funnelcmp.FunnelBuilderPage(cmp), funnelcmp.FunnelBuilder(cmp))

	return &cmp, nil
}

func formatListFunnelsResponse(ctx context.Context, response interface{}) (components.Component, error) {
	res, ok := response.([]funnel.Summary)
	if !ok {
		return &components.BaseComponent{}, fmt.Errorf("formatListFunnelsResponse: failed to cast response")
	}

	cmp := funnelcmp.FunnelsListView{
		Funnels: funnelcmp.NewFunnelSummaryView(res),
	}
	cmp.SetTemplates(funnelcmp.FunnelsListPage(cmp), funnelcmp.FunnelsList(cmp))

	return &cmp, nil
}

func formatSearchFlowsResponse(ctx context.Context, response interface{}) (components.Component, error) {
	res, ok := response.(searchFlowsResponse)
	if !ok {
		return &components.BaseComponent{}, fmt.Errorf("formatSearchFlowsResponse: failed to cast response")
	}

	cmp := funnelcmp.NewFlowSearchResultsView(res.FunnelID, res.Flows)
	cmp.SetTemplates(funnelcmp.FlowsList(cmp), funnelcmp.FlowsList(cmp))

	return &cmp, nil
}

func formatAddFlowToFunnelResponse(ctx context.Context, response interface{}) (components.Component, error) {
	res, ok := response.(funnelFlowsResponse)
	if !ok {
		return &components.BaseComponent{}, fmt.Errorf("formatAddFlowToFunnelResponse: failed to cast response")
	}

	cmp, err := formatFunnelFlowsResponse(ctx, res, true)
	if err != nil {
		return &components.BaseComponent{}, fmt.Errorf("formatAddFlowToFunnelResponse: %w", err)
	}

	return cmp, nil
}

func formatRemoveFlowFromFunnelResponse(ctx context.Context, response interface{}) (components.Component, error) {
	res, ok := response.(funnelFlowsResponse)
	if !ok {
		return &components.BaseComponent{}, fmt.Errorf("formatRemoveFlowFromFunnelResponse: failed to cast response")
	}

	cmp, err := formatFunnelFlowsResponse(ctx, res, true)
	if err != nil {
		return &components.BaseComponent{}, fmt.Errorf("formatRemoveFlowFromFunnelResponse: %w", err)
	}

	return cmp, nil
}

func formatFunnelFlowsResponse(ctx context.Context, res funnelFlowsResponse, isOOB bool) (components.Component, error) {
	flows := make([]funnelcmp.FunnelFlowView, 0)
	for _, f := range res.Flows {
		flows = append(flows, funnelcmp.FunnelFlowView{
			ID:   f.FlowID,
			Name: f.Name,
		})
	}

	cmp := funnelcmp.FunnelFormView{
		FunnelView: funnelcmp.FunnelView{
			ID: res.FunnelID,
		},
		Flows: flows,
	}
	cmp.SetUser(auth.UserFromContext(ctx))

	if len(res.Errors) > 0 {
		cmp.Errors = make([]components.ErrorView, 0)
		for _, err := range res.Errors {
			cmp.Errors = append(cmp.Errors, components.ErrorView{
				Error: err.Error(),
			})
		}
	}

	cmp.SetTemplates(funnelcmp.FunnelFlows(cmp, isOOB), funnelcmp.FunnelFlows(cmp, isOOB))

	return &cmp, nil
}
