package flowing

import (
	"context"
	"fmt"
	flowcmp "github.com/itsnoproblem/prmry/internal/components/flow"
	"github.com/itsnoproblem/prmry/internal/components/redirect"
	"github.com/itsnoproblem/prmry/internal/flow"
	"net/http"
	"regexp"
	"time"

	"github.com/pkg/errors"

	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/itsnoproblem/prmry/internal/htmx"
)

type listFlowsResponse struct {
	Summaries []flowSummary
}

type flowSummary struct {
	ID          string
	Name        string
	RuleCount   int
	LastChanged time.Time
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

func makeNewFlowFormEndpoint(svc Service) htmx.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		user, err := getAuthorizedUser(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "makeNewFlowFormEndpoint")
		}

		flows, err := svc.GetFlowsForUser(ctx, user.ID)
		if err != nil {
			return nil, errors.Wrap(err, "makeNewFlowFormEndpoint")
		}

		cmp := flowcmp.Detail{
			SupportedFields:     flow.SupportedFields(),
			SupportedConditions: flow.SupportedConditions(),
		}
		if existing := ctx.Value("view"); existing != nil {
			cmp = existing.(flowcmp.Detail)
		}

		cmp.AvailableFlowsByID = make(map[string]string)
		for _, flw := range flows {
			cmp.AvailableFlowsByID[flw.ID] = flw.Name
		}

		cmp.SetUser(auth.UserFromContext(ctx))

		fullPage := flowcmp.FlowBuilderPage(cmp)
		fragment := flowcmp.FlowBuilder(cmp)
		cmp.SetTemplates(fullPage, fragment)

		return flowBuilderResponse{
			Form: cmp,
		}, nil
	}
}

type editFlowRequest struct {
	FlowID string
	Form   *flowcmp.Detail
}

func makeEditFlowFormEndpoint(svc Service) htmx.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		user, err := getAuthorizedUser(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "makeEditFlowFormEndpoint")
		}

		req, ok := request.(editFlowRequest)
		if !ok {
			return nil, errors.Wrap(err, "makeEditFlowFormEndpoint")
		}

		var cmp flowcmp.Detail
		if existing := ctx.Value("view"); existing != nil {
			cmp = existing.(flowcmp.Detail)
		} else {
			flw, err := svc.GetFlow(ctx, req.FlowID)
			if err != nil {
				return nil, errors.Wrap(err, "makeEditFlowFormEndpoint")
			}

			cmp = flowcmp.NewDetail(flw)
		}

		flows, err := svc.GetFlowsForUser(ctx, user.ID)
		if err != nil {
			return nil, errors.Wrap(err, "makeEditFlowFormEndpoint")
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

func makeFlowBuilderAddRuleEndpoint(svc Service) htmx.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		cmp, ok := request.(flowcmp.Detail)
		if !ok {
			return nil, fmt.Errorf("makeFlowBuilderAddRuleEndpoint: failed to parse request")
		}

		cmp.Rules = append(cmp.Rules, flowcmp.RuleView{})

		return flowBuilderResponse{
			Form: cmp,
		}, nil
	}
}

type flowBuilderRemoveRuleRequest struct {
	Index int
	Form  flowcmp.Detail
}

func makeFlowBuilderRemoveRuleEndpoint(svc Service) htmx.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(flowBuilderRemoveRuleRequest)
		if !ok {
			return nil, fmt.Errorf("makeFlowBuilderRemoveRuleEndpoint: failed to parse request")
		}

		cmp := req.Form

		conditions := make([]flowcmp.RuleView, 0)
		conditions = append(conditions, cmp.Rules[:req.Index]...)
		cmp.Rules = append(conditions, cmp.Rules[req.Index+1:]...)

		return flowBuilderResponse{
			Form: cmp,
		}, nil
	}
}

func makeFlowBuilderUpdatePromptEndpoint(svc Service) htmx.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		cmp, ok := request.(flowcmp.Detail)
		if !ok {
			return nil, fmt.Errorf("makeFlowBuilderUpdatePromptEndpoint: failed to parse request")
		}

		user, err := getAuthorizedUser(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "makeFlowBuilderUpdatePromptEndpoint")
		}

		flows, err := svc.GetFlowsForUser(ctx, user.ID)
		if err != nil {
			return nil, errors.Wrap(err, "makeFlowBuilderUpdatePromptEndpoint")
		}

		cmp.AvailableFlowsByID = make(map[string]string)
		for _, flw := range flows {
			cmp.AvailableFlowsByID[flw.ID] = flw.Name
		}

		re, err := regexp.Compile("%s")
		if err != nil {
			return nil, errors.Wrap(err, "makeFlowBuilderUpdatePromptEndpoint")
		}

		matches := re.FindAllString(cmp.Prompt, -1)
		if len(matches) > len(cmp.PromptArgs) {
			for i := len(cmp.PromptArgs); i < len(matches); i++ {
				cmp.PromptArgs = append(cmp.PromptArgs, flowcmp.PromptArg{})
			}
		}

		if len(matches) < len(cmp.PromptArgs) {
			cmp.PromptArgs = cmp.PromptArgs[:len(matches)]
		}

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
