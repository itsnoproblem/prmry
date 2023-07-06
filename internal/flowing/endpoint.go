package flowing

import (
	"context"
	"fmt"
	flowcmp "github.com/itsnoproblem/prmry/internal/components/flow"
	"github.com/itsnoproblem/prmry/internal/components/redirect"
	"github.com/itsnoproblem/prmry/internal/flow"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/itsnoproblem/prmry/internal/htmx"
)

type flowSummary struct {
	ID          string
	Name        string
	RuleCount   int
	LastChanged time.Time
}

type listFlowsResponse struct {
	Summaries []flowSummary
}

func makeListFlowsEndpoint(svc Service) htmx.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		user, err := getAuthorizedUser(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "makeListFlowsEndpoint")
		}

		flows, err := svc.GetFlowsForUser(ctx, user.ID)
		if err != nil {
			return nil, errors.Wrap(err, "flowing.makeListFlowsEndpoint")
		}

		summaries := make([]flowSummary, 0)
		for _, flow := range flows {
			summaries = append(summaries, flowSummary{
				ID:          flow.ID,
				Name:        flow.Name,
				RuleCount:   len(flow.Rules),
				LastChanged: flow.UpdatedAt,
			})
		}

		return listFlowsResponse{
			Summaries: summaries,
		}, nil
	}
}

type flowBuilderResponse struct {
	Form flowcmp.Detail
}

type editFlowRequest struct {
	FlowID string
	Form   *flowcmp.Detail
}

func makeEditFlowEndpoint(svc Service) htmx.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		user, err := getAuthorizedUser(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "makeEditFlowEndpoint")
		}

		req, ok := request.(editFlowRequest)
		if !ok {
			return nil, errors.Wrap(err, "makeEditFlowEndpoint")
		}

		var cmp flowcmp.Detail
		if existing := ctx.Value("view"); existing != nil {
			cmp = existing.(flowcmp.Detail)
		} else {
			flw, err := svc.GetFlow(ctx, req.FlowID)
			if err != nil {
				return nil, errors.Wrap(err, "makeEditFlowEndpoint")
			}

			cmp = flowcmp.NewDetail(flw)
		}

		flows, err := svc.GetFlowsForUser(ctx, user.ID)
		if err != nil {
			return nil, errors.Wrap(err, "makeEditFlowEndpoint")
		}

		cmp.AvailableFlowsByID = make(map[string]string)
		for _, flw := range flows {
			cmp.AvailableFlowsByID[flw.ID] = flw.Name
		}

		cmp.SupportedFields = flow.SupportedFields()
		cmp.SupportedConditions = flow.SupportedConditions()

		return flowBuilderResponse{
			Form: cmp,
		}, nil
	}
}

type successMessageResponse struct {
	Message string
}

func makeSaveFlowEndpoint(svc Service) htmx.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_, err = getAuthorizedUser(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "makeSaveFlowEndpoint")
		}

		req, ok := request.(flowcmp.Detail)
		if !ok {
			return nil, fmt.Errorf("makeSaveFlowEndpoint: failed to parse form")
		}

		flw := req.ToFlow()

		if req.ID == "" {
			if flw.ID, err = svc.CreateFlow(ctx, flw); err != nil {
				return nil, errors.Wrap(err, "makeSaveFlowEndpoint")
			}
		} else {
			if err = svc.UpdateFlow(ctx, flw); err != nil {
				return nil, errors.Wrap(err, "makeSaveFlowEndpoint")
			}
		}

		return redirect.View{
			Location: "/flows",
			Status:   http.StatusFound,
		}, nil
	}
}

type deleteFlowRequest struct {
	FlowID string
}

func (req deleteFlowRequest) validate() error {
	if req.FlowID == "" {
		return fmt.Errorf("mising FlowID")
	}
	return nil
}

func makeDeleteFlowEndpoint(svc Service) htmx.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_, err = getAuthorizedUser(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "makeDeleteFlowEndpoint")
		}

		req, ok := request.(deleteFlowRequest)
		if !ok {
			return nil, fmt.Errorf("makeDeleteFlowEndpoint: failed to parse request")
		}

		if err := svc.DeleteFlow(ctx, req.FlowID); err != nil {
			return nil, errors.Wrap(err, "makeDeleteFlowEndpoint")
		}

		return successMessageResponse{
			Message: "Deleted flow",
		}, nil
	}
}

func getAuthorizedUser(ctx context.Context) (auth.User, error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return auth.User{}, fmt.Errorf("user is missing")
	}

	return *user, nil
}
