package inmem

import (
	"context"
	"github.com/itsnoproblem/prmry/internal/flow"
	"github.com/pkg/errors"
)

type FlowsRepo struct {
	flows     map[string]flow.Flow   // keyed by flowID
	userFlows map[string][]flow.Flow // keyed by userID
}

func NewFlowsRepo() *FlowsRepo {
	return &FlowsRepo{
		flows:     make(map[string]flow.Flow),
		userFlows: make(map[string][]flow.Flow),
	}
}

func (r *FlowsRepo) InsertFlow(ctx context.Context, flw flow.Flow) error {
	if _, ok := r.flows[flw.ID]; ok {
		return errors.New("flow already exists")
	}

	r.flows[flw.ID] = flw
	r.userFlows[flw.UserID] = append(r.userFlows[flw.UserID], flw)

	return nil
}

func (r *FlowsRepo) UpdateFlow(ctx context.Context, flw flow.Flow) error {
	if _, ok := r.flows[flw.ID]; !ok {
		return errors.New("flow does not exist")
	}

	r.flows[flw.ID] = flw

	// Update flow in userFlows map
	for i, uf := range r.userFlows[flw.UserID] {
		if uf.ID == flw.ID {
			r.userFlows[flw.UserID][i] = flw
			break
		}
	}

	return nil
}

func (r *FlowsRepo) DeleteFlow(ctx context.Context, flowID string) error {
	flw, ok := r.flows[flowID]
	if !ok {
		return errors.New("flow not found")
	}

	delete(r.flows, flowID)

	// Remove flow from userFlows
	for i, uf := range r.userFlows[flw.UserID] {
		if uf.ID == flowID {
			r.userFlows[flw.UserID] = append(r.userFlows[flw.UserID][:i], r.userFlows[flw.UserID][i+1:]...)
			break
		}
	}

	return nil
}

func (r *FlowsRepo) GetFlow(ctx context.Context, flowID string) (flow.Flow, error) {
	flw, ok := r.flows[flowID]
	if !ok {
		return flow.Flow{}, errors.New("flow not found")
	}

	return flw, nil
}

func (r *FlowsRepo) GetFlowsForUser(ctx context.Context, userID string) ([]flow.Flow, error) {
	userFlows, ok := r.userFlows[userID]
	if !ok || len(userFlows) == 0 {
		return nil, errors.New("no flows found for the given user")
	}

	flows := make([]flow.Flow, len(userFlows))
	for i, uf := range userFlows {
		flows[i] = uf
	}

	return flows, nil
}
