package funneling

import (
	"context"
	"github.com/google/uuid"
	"github.com/itsnoproblem/prmry/internal/funnel"
	"github.com/pkg/errors"
	"time"
)

type FunnelsRepository interface {
	InsertFunnel(ctx context.Context, f funnel.Funnel) error
	AddFlowsToFunnel(ctx context.Context, funnelID string, flowIDs ...string) error
	GetFunnelSummariesForUser(ctx context.Context, userID string) ([]funnel.Summary, error)
}

type service struct {
	repo FunnelsRepository
}

func NewService(repo FunnelsRepository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateFunnel(ctx context.Context, funnel funnel.Funnel) (string, error) {
	funnel.ID = uuid.NewString()
	funnel.CreatedAt = time.Now()
	funnel.UpdatedAt = funnel.CreatedAt

	if err := s.repo.InsertFunnel(ctx, funnel); err != nil {
		return "", errors.Wrap(err, "funneling.CreateFunnel")
	}

	return funnel.ID, nil
}

func (s *service) ListFunnels(ctx context.Context, userID string) ([]funnel.Summary, error) {
	funnels, err := s.repo.GetFunnelSummariesForUser(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "funneling.ListFunnels")
	}

	return funnels, nil
}

func (s *service) AddFlowsToFunnel(ctx context.Context, funnelID string, flowIDs ...string) error {
	if err := s.repo.AddFlowsToFunnel(ctx, funnelID, flowIDs...); err != nil {
		return errors.Wrap(err, "funneling.AddFlowsToFunnel")
	}

	return nil
}
