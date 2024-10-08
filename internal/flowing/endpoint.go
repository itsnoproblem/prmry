package flowing

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/itsnoproblem/prmry/internal/funnel"
	"github.com/pkg/errors"

	"github.com/itsnoproblem/prmry/internal/auth"
	flowcmp "github.com/itsnoproblem/prmry/internal/components/flow"
	"github.com/itsnoproblem/prmry/internal/components/redirect"
	"github.com/itsnoproblem/prmry/internal/flow"
	internalhttp "github.com/itsnoproblem/prmry/internal/http"
)

const (
	ContextKeyFlow = "view"
)

type Service interface {
	CreateFlow(ctx context.Context, flw flow.Flow) (ID string, err error)
	UpdateFlow(ctx context.Context, flw flow.Flow) error
	DeleteFlow(ctx context.Context, flowID string) error
	GetFlow(ctx context.Context, flowID string) (flow.Flow, error)
	GetFlowsForUser(ctx context.Context, userID string) ([]flow.Flow, error)
	GetFunnelsForFlow(ctx context.Context, flowID string) ([]funnel.Funnel, error)
	APIURL() string
}

func makeListFlowsEndpoint(svc Service) internalhttp.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		user, err := getAuthorizedUser(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "makeListFlowsEndpoint")
		}

		flows, err := svc.GetFlowsForUser(ctx, user.ID)
		if err != nil {
			return nil, errors.Wrap(err, "flowing.makeListFlowsEndpoint")
		}

		return flows, nil
	}
}

type flowBuilderResponse struct {
	Form flowcmp.Detail
}

func makeFlowBuilderEndpoint(svc Service) internalhttp.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(flowcmp.Detail)
		if !ok {
			req = flowcmp.NewDetail(flow.Flow{})
		}

		user, err := getAuthorizedUser(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "makeFlowBuilderEndpoint")
		}

		flows, err := svc.GetFlowsForUser(ctx, user.ID)
		if err != nil {
			return nil, errors.Wrap(err, "makeFlowBuilderEndpoint")
		}

		cmp := req
		if existing := ctx.Value(ContextKeyFlow); existing != nil {
			cmp = existing.(flowcmp.Detail)
		}

		cmp.SetAvalableFlows(flows)
		cmp.SetUser(&user)
		cmp.SelectedTab = req.SelectedTab
		cmp.FlowURL = svc.APIURL() + "/flows/" + req.ID

		associatedFunnels, err := svc.GetFunnelsForFlow(ctx, req.ID)
		if err != nil {
			return nil, errors.Wrap(err, "makeFlowBuilderEndpoint")
		}

		cmp.Funnels = make([]flowcmp.Funnel, 0)
		for _, fnl := range associatedFunnels {
			cmp.Funnels = append(cmp.Funnels, flowcmp.Funnel{
				Name: fnl.Name,
				URL:  svc.APIURL() + "/funnels/" + fnl.Path,
			})
		}

		fullPage := flowcmp.FlowBuilderPage(cmp)
		fragment := flowcmp.FlowBuilder(cmp)
		cmp.SetTemplates(fullPage, fragment)

		return flowBuilderResponse{
			Form: cmp,
		}, nil
	}
}

type editFlowRequest struct {
	FlowID      string
	Form        *flowcmp.Detail
	SelectedTab string
}

func makeEditFlowFormEndpoint(svc Service) internalhttp.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		user, err := getAuthorizedUser(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "makeEditFlowFormEndpoint")
		}

		req, ok := request.(editFlowRequest)
		if !ok {
			return nil, errors.Wrap(err, "makeEditFlowFormEndpoint")
		}

		cmp := flowcmp.NewDetail(flow.Flow{})
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

		cmp.SetAvalableFlows(flows)
		cmp.SetUser(&user)

		cmp.SupportedFields = flow.SupportedFields()
		cmp.SupportedConditions = flow.SupportedConditions()
		cmp.FlowURL = svc.APIURL() + "/flows/" + req.FlowID
		cmp.SelectedTab = req.SelectedTab

		associatedFunnels, err := svc.GetFunnelsForFlow(ctx, req.FlowID)
		if err != nil {
			return nil, errors.Wrap(err, "makeEditFlowFormEndpoint")
		}

		cmp.Funnels = make([]flowcmp.Funnel, 0)
		for _, fnl := range associatedFunnels {
			cmp.Funnels = append(cmp.Funnels, flowcmp.Funnel{
				Name: fnl.Name,
				URL:  svc.APIURL() + "/funnels/" + fnl.Path,
			})
		}

		return flowBuilderResponse{
			Form: cmp,
		}, nil
	}
}

func makeFlowBuilderAddInputEndpoint(svc Service) internalhttp.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		cmp, ok := request.(flowcmp.Detail)
		if !ok {
			return nil, fmt.Errorf("makeFlowBuilderAddRuleEndpoint: failed to parse request")
		}

		if cmp.InputParams == nil {
			cmp.InputParams = make([]flowcmp.InputParam, 0)
		}
		cmp.InputParams = append(cmp.InputParams, flowcmp.InputParam{})
		return flowBuilderResponse{
			Form: cmp,
		}, nil
	}
}

func makeFlowBuilderAddRuleEndpoint(svc Service) internalhttp.HandlerFunc {
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

type flowBuilderRemoveItemRequest struct {
	Index int
	Form  flowcmp.Detail
}

func makeFlowBuilderRemoveInputEndpoint(svc Service) internalhttp.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(flowBuilderRemoveItemRequest)
		if !ok {
			return nil, fmt.Errorf("makeFlowBuilderRemoveInputEndpoint: failed to parse request")
		}

		cmp := req.Form

		if len(cmp.Rules) > 0 {
			revisedInputs := make([]flowcmp.InputParam, 0)
			revisedInputs = append(revisedInputs, cmp.InputParams[:req.Index]...)
			cmp.InputParams = append(revisedInputs, cmp.InputParams[req.Index+1:]...)
		}

		return flowBuilderResponse{
			Form: cmp,
		}, nil
	}
}

func makeFlowBuilderRemoveRuleEndpoint(svc Service) internalhttp.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(flowBuilderRemoveItemRequest)
		if !ok {
			return nil, fmt.Errorf("makeFlowBuilderRemoveRuleEndpoint: failed to parse request")
		}

		cmp := req.Form

		if len(cmp.Rules) > 0 {
			revisedRules := make([]flowcmp.RuleView, 0)
			revisedRules = append(revisedRules, cmp.Rules[:req.Index]...)
			cmp.Rules = append(revisedRules, cmp.Rules[req.Index+1:]...)
		}

		return flowBuilderResponse{
			Form: cmp,
		}, nil
	}
}

func makeFlowBuilderUpdatePromptEndpoint(svc Service) internalhttp.HandlerFunc {
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

func makeSaveFlowEndpoint(svc Service) internalhttp.HandlerFunc {
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
	FlowID      string
	SelectedTab string
}

func (req deleteFlowRequest) validate() error {
	if req.FlowID == "" {
		return fmt.Errorf("mising FlowID")
	}
	return nil
}

func makeDeleteFlowEndpoint(svc Service) internalhttp.HandlerFunc {
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
