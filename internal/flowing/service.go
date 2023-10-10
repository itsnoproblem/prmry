package flowing

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/itsnoproblem/prmry/internal/flow"
)

type FlowsRepo interface {
	InsertFlow(ctx context.Context, flw flow.Flow) error
	UpdateFlow(ctx context.Context, flw flow.Flow) error
	DeleteFlow(ctx context.Context, flowID string) error
	GetFlow(ctx context.Context, flowID string) (flow.Flow, error)
	GetFlowsForUser(ctx context.Context, userID string) ([]flow.Flow, error)
}

type service struct {
	flowsRepo FlowsRepo
}

func NewService(repo FlowsRepo) *service {
	return &service{
		flowsRepo: repo,
	}
}

func (s *service) GetFlow(ctx context.Context, flowID string) (flow.Flow, error) {
	flw, err := s.flowsRepo.GetFlow(ctx, flowID)
	if err != nil {
		return flow.Flow{}, errors.Wrap(err, "flowing.GetFlow")
	}

	return flw, nil
}

func (s *service) GetFlowsForUser(ctx context.Context, userID string) ([]flow.Flow, error) {
	flows, err := s.flowsRepo.GetFlowsForUser(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "flowing.GetFlowsForUser")
	}

	return flows, nil
}

func (s *service) CreateFlow(ctx context.Context, flw flow.Flow) (ID string, err error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return "", errors.Wrap(err, "flowing.SaveFlow: user is nil")
	}

	flw.ID = uuid.NewString()
	flw.UserID = user.ID
	flw.CreatedAt = time.Now()
	flw.UpdatedAt = time.Now()

	if err = s.flowsRepo.InsertFlow(ctx, flw); err != nil {
		return "", errors.Wrap(err, "flowing.SaveFlow")
	}

	return flw.ID, nil
}

func (s *service) UpdateFlow(ctx context.Context, flw flow.Flow) error {
	flw.UpdatedAt = time.Now()
	if err := s.flowsRepo.UpdateFlow(ctx, flw); err != nil {
		return errors.Wrap(err, "flowing.UpdateFlow")
	}

	return nil
}

func (s *service) DeleteFlow(ctx context.Context, flowID string) error {
	if err := s.flowsRepo.DeleteFlow(ctx, flowID); err != nil {
		return errors.Wrap(err, "flowing.DeleteFlow")
	}

	return nil
}
