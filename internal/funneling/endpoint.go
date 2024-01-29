package funneling

import (
	"context"
	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/itsnoproblem/prmry/internal/components/redirect"
	"github.com/itsnoproblem/prmry/internal/flow"
	"github.com/itsnoproblem/prmry/internal/funnel"
	internalhttp "github.com/itsnoproblem/prmry/internal/http"
	"github.com/pkg/errors"
	"net/http"
)

type Service interface {
	CreateFunnel(ctx context.Context, funnel funnel.Funnel) (string, error)
	UpdateFunnel(ctx context.Context, funnel funnel.Funnel) error
	DeleteFunnel(ctx context.Context, funnelID string) error
	ListFunnels(ctx context.Context, userID string) ([]funnel.Summary, error)
	AddFlowsToFunnel(ctx context.Context, funnelID string, flowIDs ...string) error
	RemoveFlowsFromFunnel(ctx context.Context, funnelID string, flowIDs ...string) error
	SearchFlows(ctx context.Context, userID, search string) ([]flow.Flow, error)
	GetFunnelWithFlows(ctx context.Context, funnelID string) (funnel.WithFlows, error)
}

type createFunnelRequest struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func (r createFunnelRequest) Validate() error {
	if r.Path == "" {
		return errors.New("funneling.createFunnelRequest: path is required")
	}

	return nil
}

func makeCreateFunnelEndpoint(svc Service) internalhttp.HandlerFunc {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		user := auth.UserFromContext(ctx)
		if user == nil {
			return nil, errors.New("funneling.makeCreateFunnelEndpoint: user not found")
		}

		req, ok := request.(createFunnelRequest)
		if !ok {
			return nil, errors.New("funneling.makeCreateFunnelEndpoint: invalid request")
		}

		if err := req.Validate(); err != nil {
			return nil, errors.Wrap(err, "funneling.makeCreateFunnelEndpoint")
		}

		newFunnel := funnel.Funnel{
			UserID: user.ID,
			Name:   req.Name,
			Path:   normalizePath(req.Path),
		}

		var err error
		newFunnel.ID, err = svc.CreateFunnel(ctx, newFunnel)
		if err != nil {
			return nil, errors.Wrap(err, "funneling.makeCreateFunnelEndpoint")
		}

		return redirect.View{
			Status:   http.StatusTemporaryRedirect,
			Location: "/funnels/" + newFunnel.ID,
		}, nil
	}
}

type funnelIDRequest struct {
	ID string
}

func makeEditFunnelEndpoint(svc Service) internalhttp.HandlerFunc {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		user := auth.UserFromContext(ctx)
		if user == nil {
			return funnel.Funnel{}, errors.New("funneling.makeEditFunnelEndpoint: user not found")
		}

		req, ok := request.(funnelIDRequest)
		if !ok {
			return funnel.Funnel{}, errors.New("funneling.makeEditFunnelEndpoint: invalid request")
		}

		fnl, err := svc.GetFunnelWithFlows(ctx, req.ID)
		if err != nil {
			return funnel.Funnel{}, errors.Wrap(err, "funneling.makeEditFunnelEndpoint")
		}

		return fnl, nil
	}
}

type funnelBuilderRequest struct {
	ID   string
	Name string
	Path string
}

func makeFunnelBuilderEndpoint(svc Service) internalhttp.HandlerFunc {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		user := auth.UserFromContext(ctx)
		if user == nil {
			return nil, errors.New("funneling.makeFunnelBuilderEndpoint: user not found")
		}

		req, ok := request.(funnelBuilderRequest)
		if !ok {
			return nil, errors.New("funneling.makeFunnelBuilderEndpoint: invalid request")
		}

		return funnel.WithFlows{
			Funnel: funnel.Funnel{
				UserID: user.ID,
				ID:     req.ID,
				Name:   req.Name,
				Path:   req.Path,
			},
		}, nil
	}
}

func (r funnelBuilderRequest) Validate() error {
	if r.Name == "" {
		return errors.New("funneling.funnelBuilderRequest: name is required")
	}

	if r.Path == "" {
		return errors.New("funneling.funnelBuilderRequest: path is required")
	}

	return nil
}

func makeSaveFunnelEndpoint(svc Service) internalhttp.HandlerFunc {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		user := auth.UserFromContext(ctx)
		if user == nil {
			return funnel.Funnel{}, errors.New("funneling.makeSaveFunnelEndpoint: user not found")
		}

		req, ok := request.(funnelBuilderRequest)
		if !ok {
			return funnel.Funnel{}, errors.New("funneling.makeSaveFunnelEndpoint: invalid request")
		}

		var err error

		funnelWithFlows := funnel.WithFlows{
			Funnel: funnel.Funnel{
				UserID: user.ID,
				ID:     req.ID,
				Name:   req.Name,
				Path:   req.Path,
			},
			Flows: make([]flow.Flow, 0),
		}

		// @TODO(marty): client side validation first
		//if err := req.Validate(); err != nil {
		//	return fnl, errors.Wrap(err, "funneling.makeSaveFunnelEndpoint")
		//}

		if req.ID == "" {
			funnelWithFlows.ID, err = svc.CreateFunnel(ctx, funnelWithFlows.Funnel)
		} else {
			funnelWithFlows, err = svc.GetFunnelWithFlows(ctx, req.ID)
			if err != nil {
				return funnel.WithFlows{}, errors.Wrap(err, "funneling.makeSaveFunnelEndpoint")
			}

			funnelWithFlows.Name = req.Name
			funnelWithFlows.Path = req.Path

			err = svc.UpdateFunnel(ctx, funnelWithFlows.Funnel)
		}

		if err != nil {
			return funnel.WithFlows{}, errors.Wrap(err, "funneling.makeSaveFunnelEndpoint")
		}

		return funnelWithFlows, nil
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

type searchFlowsRequest struct {
	FunnelID string
	Search   string
}

type searchFlowsResponse struct {
	FunnelID string
	Flows    []flow.Flow
}

func makeSearchFlowsEndpoint(svc Service) internalhttp.HandlerFunc {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		user := auth.UserFromContext(ctx)
		if user == nil {
			return nil, errors.New("funneling.makeSearchFlowsEndpoint: user not found")
		}

		req, ok := request.(searchFlowsRequest)
		if !ok {
			return nil, errors.New("funneling.makeSearchFlowsEndpoint: invalid request")
		}

		flows, err := svc.SearchFlows(ctx, user.ID, req.Search)
		if err != nil {
			return nil, errors.Wrap(err, "funneling.makeSearchFlowsEndpoint")
		}

		return searchFlowsResponse{
			FunnelID: req.FunnelID,
			Flows:    flows,
		}, nil
	}
}

type addFlowToFunnelRequest struct {
	FunnelID string
	FlowID   string
}

func (r addFlowToFunnelRequest) Validate() error {
	if r.FunnelID == "" {
		return errors.New("funneling.addFlowToFunnelRequest: funnel id is required")
	}

	if r.FlowID == "" {
		return errors.New("funneling.addFlowToFunnelRequest: flow id is required")
	}

	return nil
}

type funnelFlowResponse struct {
	FlowID string
	Name   string
}

type funnelFlowsResponse struct {
	FunnelID string
	Flows    []funnelFlowResponse
	Errors   []error
}

func newFunnelFlowsResponse(funnelID string, flows []flow.Flow) funnelFlowsResponse {
	res := funnelFlowsResponse{
		FunnelID: funnelID,
		Flows:    make([]funnelFlowResponse, 0),
	}

	for _, f := range flows {
		res.Flows = append(res.Flows, funnelFlowResponse{
			FlowID: f.ID,
			Name:   f.Name,
		})
	}

	return res
}

func makeAddFlowToFunnelEndpoint(svc Service) internalhttp.HandlerFunc {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		user := auth.UserFromContext(ctx)
		if user == nil {
			return nil, errors.New("funneling.makeAddFlowToFunnelEndpoint: user not found")
		}

		req, ok := request.(addFlowToFunnelRequest)
		if !ok {
			return nil, errors.New("funneling.makeAddFlowToFunnelEndpoint: invalid request")
		}

		if err := req.Validate(); err != nil {
			return nil, errors.Wrap(err, "funneling.makeAddFlowToFunnelEndpoint")
		}

		var allErr []error
		if err := svc.AddFlowsToFunnel(ctx, req.FunnelID, req.FlowID); err != nil {
			allErr = append(allErr, errors.Wrap(err, "funneling.makeAddFlowToFunnelEndpoint"))
		}

		fnl, err := svc.GetFunnelWithFlows(ctx, req.FunnelID)
		if err != nil {
			allErr = append(allErr, errors.Wrap(err, "funneling.makeAddFlowToFunnelEndpoint"))
		}

		res := newFunnelFlowsResponse(fnl.ID, fnl.Flows)
		res.Errors = allErr

		return res, nil
	}
}

type removeFlowFromFunnelRequest struct {
	FunnelID string
	FlowID   string
}

func makeRemoveFlowFromFunnelEndpoint(svc Service) internalhttp.HandlerFunc {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		user := auth.UserFromContext(ctx)
		if user == nil {
			return nil, errors.New("funneling.makeRemoveFlowFromFunnelEndpoint: user not found")
		}

		req, ok := request.(removeFlowFromFunnelRequest)
		if !ok {
			return nil, errors.New("funneling.makeRemoveFlowFromFunnelEndpoint: invalid request")
		}

		if err := svc.RemoveFlowsFromFunnel(ctx, req.FunnelID, req.FlowID); err != nil {
			return nil, errors.Wrap(err, "funneling.makeRemoveFlowFromFunnelEndpoint")
		}

		fnl, err := svc.GetFunnelWithFlows(ctx, req.FunnelID)
		if err != nil {
			return nil, errors.Wrap(err, "funneling.makeRemoveFlowFromFunnelEndpoint")
		}

		return newFunnelFlowsResponse(fnl.ID, fnl.Flows), nil
	}
}

func makeDeleteFunnelEndpoint(svc Service) internalhttp.HandlerFunc {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(funnelIDRequest)
		if !ok {
			return nil, errors.New("funneling.makeDeleteFunnelEndpoint: invalid request")
		}

		if err := svc.DeleteFunnel(ctx, req.ID); err != nil {
			return nil, errors.Wrap(err, "funneling.makeDeleteFunnelEndpoint")
		}

		return redirect.View{
			Status:   http.StatusTemporaryRedirect,
			Location: "/funnels",
		}, nil
	}
}
