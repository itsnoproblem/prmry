package flowing

import (
	"context"
	"github.com/google/uuid"
	"github.com/itsnoproblem/prmry/pkg/auth"
	"github.com/pkg/errors"
	"time"

	"github.com/itsnoproblem/prmry/pkg/flow"
)

type FlowsRepo interface {
	InsertFlow(ctx context.Context, flw flow.Flow) error
}

type service struct {
	flowsRepo FlowsRepo
}

func NewService(repo FlowsRepo) *service {
	return &service{
		flowsRepo: repo,
	}
}

func (s *service) GetFlowsForUser(ctx context.Context, userID string) ([]flow.Flow, error) {
	return nil, nil
}

func (s *service) CreateFlow(ctx context.Context, flw flow.Flow) (ID string, err error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return "", errors.Wrap(err, "flowing.CreateFlow: user is nil")
	}

	flw.ID = uuid.NewString()
	flw.UserID = user.ID
	flw.CreatedAt = time.Now()
	flw.UpdatedAt = time.Now()

	if err = s.flowsRepo.InsertFlow(ctx, flw); err != nil {
		return "", errors.Wrap(err, "flowing.CreateFlow")
	}

	return flw.ID, nil
}

func (s service) UpdateFlow(ctx context.Context, flw flow.Flow) error {
	return nil
}
