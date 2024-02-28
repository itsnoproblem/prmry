package inmem

import (
    "context"

    "github.com/itsnoproblem/prmry/internal/funnel"
)

type funnelsRepo struct {
    funnelsByIDs map[string][]funnel.Funnel
}

func NewFunnelsRepo() *funnelsRepo {
    return &funnelsRepo{
        funnelsByIDs: make(map[string][]funnel.Funnel),
    }
}

func (r *funnelsRepo) GetFunnelsForFlow(ctx context.Context, flowID string) ([]funnel.Funnel, error) {
    return r.funnelsByIDs[flowID], nil
}
