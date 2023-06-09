package flowing

import (
	"context"

	"github.com/itsnoproblem/prmry/pkg/flow"
)

type service struct{}

func NewService() *service {
	return &service{}
}

func (s *service) GetFlowsForUser(ctx context.Context, userID string) ([]flow.Flow, error) {
	return nil, nil
}
