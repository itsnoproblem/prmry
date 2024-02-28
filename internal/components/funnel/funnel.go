package funnel

import (
	"fmt"
	"github.com/itsnoproblem/prmry/internal/components"
	"github.com/itsnoproblem/prmry/internal/flow"
	"github.com/itsnoproblem/prmry/internal/funnel"
)

type FunnelView struct {
	ID   string
	Name string
	Path string
}

type FunnelSummaryView struct {
	FunnelView
	FlowCount string
}

type FunnelFormView struct {
	components.BaseComponent
	FunnelView
	Flows []FunnelFlowView
}

type FunnelFlowView struct {
	components.BaseComponent
	ID   string
	Name string
}

type FunnelsListView struct {
	components.BaseComponent
	Funnels []FunnelSummaryView
}

func NewFunnelSummaryView(summaries []funnel.Summary) []FunnelSummaryView {
	view := make([]FunnelSummaryView, 0)
	for _, s := range summaries {
		view = append(view, FunnelSummaryView{
			FunnelView: FunnelView{
				ID:   s.ID,
				Name: s.Name,
				Path: s.Path,
			},
			FlowCount: fmt.Sprintf("%d", s.FlowCount),
		})
	}

	return view
}

type FlowSearchResultsView struct {
	components.BaseComponent
	FunnelID string
	Flows    []FlowSearchResultView
}

type FlowSearchResultView struct {
	components.BaseComponent
	ID   string
	Name string
}

func NewFlowSearchResultsView(funnelID string, flows []flow.Flow) FlowSearchResultsView {
	view := FlowSearchResultsView{
		FunnelID: funnelID,
	}

	for _, f := range flows {
		view.Flows = append(view.Flows, FlowSearchResultView{
			ID:   f.ID,
			Name: f.Name,
		})
	}

	return view
}
