package funneling

import (
	"context"
	"github.com/pkg/errors"

	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/itsnoproblem/prmry/internal/funnel"
	internalhttp "github.com/itsnoproblem/prmry/internal/http"
)

type Service interface {
	CreateFunnel(ctx context.Context, funnel funnel.Funnel) (string, error)
	ListFunnels(ctx context.Context, userID string) ([]funnel.Summary, error)
	AddFlowsToFunnel(ctx context.Context, funnelID string, flowIDs ...string) error
}

func makeFunnelBuilderEndpoint(svc Service) internalhttp.HandlerFunc {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		user := auth.UserFromContext(ctx)
		if user == nil {
			return nil, errors.New("funneling.makeFunnelBuilderEndpoint: user not found")
		}

		return funnel.Funnel{}, nil
	}
}

type CreateFunnelRequest struct {
	Name string
	Path string
}

type CreateFunnelResponse struct {
	ID string `json:"id"`
}

func makeCreateFunnelEndpoint(svc Service) internalhttp.HandlerFunc {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		user := auth.UserFromContext(ctx)
		if user == nil {
			return nil, errors.New("funneling.makeCreateFunnelEndpoint: user not found")
		}

		req, ok := request.(CreateFunnelRequest)
		if !ok {
			return nil, errors.New("funneling.makeCreateFunnelEndpoint: invalid request")
		}

		funnel := funnel.Funnel{
			UserID: user.ID,
			Name:   req.Name,
			Path:   req.Path,
		}

		var err error
		funnel.ID, err = svc.CreateFunnel(ctx, funnel)
		if err != nil {
			return nil, errors.Wrap(err, "funneling.makeCreateFunnelEndpoint")
		}

		return funnel, nil
	}
}

func makeListFunnelsEndpoint(svc Service) internalhttp.HandlerFunc {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		user := auth.UserFromContext(ctx)
		if user == nil {
			return nil, errors.New("funneling.makeListFunnelsEndpoint: user not found")
		}

		funnels, err := svc.ListFunnels(ctx, user.ID)
		if err != nil {
			return nil, errors.Wrap(err, "funneling.makeListFunnelsEndpoint")
		}

		return funnels, nil
	}
}
