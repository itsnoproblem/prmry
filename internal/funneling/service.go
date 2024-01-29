package funneling

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/itsnoproblem/prmry/internal/flow"
	"github.com/itsnoproblem/prmry/internal/funnel"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type FunnelsRepository interface {
	InsertFunnel(ctx context.Context, f funnel.Funnel) error
	UpdateFunnel(ctx context.Context, f funnel.Funnel) error
	DeleteFunnel(ctx context.Context, funnelID string) error
	ExistsFunnelFlow(ctx context.Context, funnelID, flowID string) (bool, error)
	AddFlowsToFunnel(ctx context.Context, funnelID string, flowIDs ...string) error
	RemoveFlowsFromFunnel(ctx context.Context, funnelID string, flowIDs ...string) error
	GetFunnel(ctx context.Context, id string) (funnel.Funnel, error)
	GetFunnelSummariesForUser(ctx context.Context, userID string) ([]funnel.Summary, error)
}

type FlowsRepository interface {
	SearchFlows(ctx context.Context, userID, search string) ([]flow.Flow, error)
	GetFlowsForFunnel(ctx context.Context, funnelID string) ([]flow.Flow, error)
}

type service struct {
	repo      FunnelsRepository
	flowsRepo FlowsRepository
}

func NewService(r FunnelsRepository, f FlowsRepository) *service {
	return &service{
		repo:      r,
		flowsRepo: f,
	}
}

func (s *service) SearchFlows(ctx context.Context, userID, search string) ([]flow.Flow, error) {
	flows, err := s.flowsRepo.SearchFlows(ctx, userID, search)
	if err != nil {
		return nil, errors.Wrap(err, "funneling.SearchFlows")
	}

	return flows, nil
}

func (s *service) CreateFunnel(ctx context.Context, fnl funnel.Funnel) (string, error) {
	fnl.ID = uuid.NewString()
	fnl.CreatedAt = time.Now()
	fnl.UpdatedAt = fnl.CreatedAt
	fnl.Path = normalizePath(fnl.Path)

	if err := s.repo.InsertFunnel(ctx, fnl); err != nil {
		return "", errors.Wrap(err, "funneling.CreateFunnel")
	}

	return fnl.ID, nil
}

func (s *service) UpdateFunnel(ctx context.Context, fnl funnel.Funnel) error {
	existing, err := s.repo.GetFunnel(ctx, fnl.ID)
	if err != nil {
		return errors.Wrap(err, "funneling.UpdateFunnel")
	}

	fnl.UpdatedAt = time.Now()
	fnl.CreatedAt = existing.CreatedAt
	fnl.UserID = existing.UserID
	fnl.Path = normalizePath(fnl.Path)

	if err := s.repo.UpdateFunnel(ctx, fnl); err != nil {
		return errors.Wrap(err, "funneling.UpdateFunnel")
	}

	return nil
}

func (s *service) DeleteFunnel(ctx context.Context, funnelID string) error {
	if err := s.repo.DeleteFunnel(ctx, funnelID); err != nil {
		return errors.Wrap(err, "funneling.DeleteFunnel")
	}

	return nil
}

func (s *service) ListFunnels(ctx context.Context, userID string) ([]funnel.Summary, error) {
	funnels, err := s.repo.GetFunnelSummariesForUser(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "funneling.ListFunnels")
	}

	return funnels, nil
}

func (s *service) AddFlowsToFunnel(ctx context.Context, funnelID string, flowIDs ...string) error {
	for _, flowID := range flowIDs {
		exists, err := s.repo.ExistsFunnelFlow(ctx, funnelID, flowID)
		if err != nil {
			return errors.Wrap(err, "funneling.AddFlowsToFunnel")
		}

		if exists {
			return fmt.Errorf("funneling.AddFlowsToFunnel: flow %s already exists in funnel", flowID)
		}
	}

	if err := s.repo.AddFlowsToFunnel(ctx, funnelID, flowIDs...); err != nil {
		return errors.Wrap(err, "funneling.AddFlowsToFunnel")
	}

	return nil
}

func (s *service) RemoveFlowsFromFunnel(ctx context.Context, funnelID string, flowIDs ...string) error {
	if err := s.repo.RemoveFlowsFromFunnel(ctx, funnelID, flowIDs...); err != nil {
		return errors.Wrap(err, "funneling.RemoveFlowsFromFunnel")
	}

	return nil
}

func (s *service) GetFunnelWithFlows(ctx context.Context, funnelID string) (funnel.WithFlows, error) {
	fnl, err := s.repo.GetFunnel(ctx, funnelID)
	if err != nil {
		return funnel.WithFlows{}, errors.Wrap(err, "funneling.GetFunnelWithFlows")
	}

	flows, err := s.flowsRepo.GetFlowsForFunnel(ctx, funnelID)
	if err != nil {
		return funnel.WithFlows{}, errors.Wrap(err, "funneling.GetFunnelWithFlows")
	}

	return funnel.WithFlows{
		Funnel: fnl,
		Flows:  flows,
	}, nil
}

func normalizePath(path string) string {
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")
	return path
}
