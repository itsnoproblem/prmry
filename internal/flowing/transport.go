package flowing

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/a-h/templ"
	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/itsnoproblem/prmry/internal/components"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"

	flowcmp "github.com/itsnoproblem/prmry/internal/components/flow"
	"github.com/itsnoproblem/prmry/internal/flow"
	"github.com/itsnoproblem/prmry/internal/htmx"
)

type Service interface {
	CreateFlow(ctx context.Context, flw flow.Flow) (ID string, err error)
	UpdateFlow(ctx context.Context, flw flow.Flow) error
	DeleteFlow(ctx context.Context, flowID string) error
	GetFlow(ctx context.Context, flowID string) (flow.Flow, error)
	GetFlowsForUser(ctx context.Context, userID string) ([]flow.Flow, error)
}

type Renderer interface {
	RenderError(w http.ResponseWriter, r *http.Request, err error)
	Render(w http.ResponseWriter, r *http.Request, cmp components.Component) error
	RenderTemplComponent(w http.ResponseWriter, r *http.Request, fullPage, fragment templ.Component) error
	Unauthorized(w http.ResponseWriter, r *http.Request)
}

func RouteHandler(svc Service, renderer Renderer) func(chi.Router) {
	listFlowsEndpoint := htmx.NewEndpoint(
		makeListFlowsEndpoint(svc),
		decodeEmptyRequest,
		formatFlowSummaries,
		auth.Required,
	)

	editFlowEndpoint := htmx.NewEndpoint(
		makeEditFlowEndpoint(svc),
		decodeEditFormRequest,
		formatFlowBuilderResponse,
		auth.Required,
	)

	saveFlowEndpoint := htmx.NewEndpoint(
		makeSaveFlowEndpoint(svc),
		decodeFlowBuilderRequest,
		formatRedirectResponse,
		auth.Required,
	)

	deleteFlowEndpoint := htmx.NewEndpoint(
		makeDeleteFlowEndpoint(svc),
		decodeDeleteFlowRequest,
		formatRedirectResponse,
		auth.Required,
	)

	return func(r chi.Router) {
		r.Post("/", htmx.MakeHandler(saveFlowEndpoint, renderer))
		r.Get("/", htmx.MakeHandler(listFlowsEndpoint, renderer))
		r.Get("/{flowID}/edit", htmx.MakeHandler(editFlowEndpoint, renderer))
		r.Delete("/{flowID}", htmx.MakeHandler(deleteFlowEndpoint, renderer))
	}
}

func decodeEmptyRequest(request *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeDeleteFlowRequest(request *http.Request) (interface{}, error) {
	return deleteFlowRequest{
		FlowID: chi.URLParam(request, "flowID"),
	}, nil
}

func decodeEditFormRequest(request *http.Request) (interface{}, error) {
	return editFlowRequest{
		FlowID: chi.URLParam(request, "flowID"),
	}, nil
}

type flowBuilderFormRequest struct {
	ID              string      `json:"id"`
	Name            string      `json:"name"`
	RequireAll      string      `json:"requireAll"`
	FieldNames      interface{} `json:"fieldName"`
	SelectedFlows   interface{} `json:"selectedFlows"`
	ConditionTypes  interface{} `json:"condition"`
	ConditionValues interface{} `json:"value"`
	RuleFlows       interface{} `json:"ruleConditionFlows"`
	Prompt          string      `json:"prompt"`
	PromptArgs      interface{} `json:"promptArgs"`
	PromptArgFlows  interface{} `json:"promptArgFlows"`
	AvailableFlows  interface{} `json:"availableFlows"`
}

func decodeFlowBuilderRequest(r *http.Request) (interface{}, error) {
	if r == nil {
		return flowcmp.Detail{}, fmt.Errorf("readForm: request was null")
	}

	var req flowBuilderFormRequest
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "decodeFlowBuilderRequest")
	}

	if err = json.Unmarshal(body, &req); err != nil {
		return nil, errors.Wrap(err, "decodeFlowBuilderRequest")
	}

	fieldNames := make([]string, 0)
	selectedFlows := make([]string, 0)
	conditionTypes := make([]string, 0)
	conditionValues := make([]string, 0)
	parsedPromptArgs := make([]string, 0)
	promptArgFlows := make([]string, 0)

	fieldNames, err = stringOrSlice(req.FieldNames)
	if err != nil {
		return nil, errors.Wrap(err, "decodeFlowBuilderRequest: fieldNames")
	}

	selectedFlows, err = stringOrSlice(req.SelectedFlows)
	if err != nil {
		return nil, errors.Wrap(err, "decodeFlowBuilderRequest: selectedFlows")
	}

	conditionTypes, err = stringOrSlice(req.ConditionTypes)
	if err != nil {
		return nil, errors.Wrap(err, "decodeFlowBuilderRequest: conditionTypes")
	}

	conditionValues, err = stringOrSlice(req.ConditionValues)
	if err != nil {
		return nil, errors.Wrap(err, "decodeFlowBuilderRequest: conditionValues")
	}

	parsedPromptArgs, err = stringOrSlice(req.PromptArgs)
	if err != nil {
		return nil, errors.Wrap(err, "decodeFlowBuilderRequest: parsedPromptArgs")
	}

	if req.PromptArgFlows != nil {
		promptArgFlows, err = stringOrSlice(req.PromptArgFlows)
		if err != nil {
			return nil, errors.Wrap(err, "decodeFlowBuilderRequest: promptArgFlows")
		}
	}

	if len(fieldNames) != len(conditionTypes) || len(fieldNames) != len(conditionValues) {
		return flowcmp.Detail{}, fmt.Errorf("readForm: condition fields mismatch")
	}

	promptArgs := make([]flowcmp.PromptArg, 0)
	flowIndex := 0
	for _, arg := range parsedPromptArgs {
		pargs := flowcmp.PromptArg{
			Source: flow.SourceType(arg),
		}

		if flow.SourceType(arg) == flow.FieldSourceFlow && len(promptArgFlows) > flowIndex {
			pargs.Value = promptArgFlows[flowIndex]
		}

		promptArgs = append(promptArgs, pargs)
	}

	form := flowcmp.Detail{
		ID:                  req.ID,
		Name:                req.Name,
		Prompt:              req.Prompt,
		PromptArgs:          promptArgs,
		SupportedFields:     flow.SupportedFields(),
		SupportedConditions: flow.SupportedConditions(),
	}

	flowIndex = 0
	form.Rules = make([]flowcmp.RuleView, 0)
	for i, name := range fieldNames {
		fieldValue := ""
		if name == flow.FieldSourceFlow.String() && len(selectedFlows) > flowIndex {
			fieldValue = selectedFlows[flowIndex]
			flowIndex++
		}

		form.Rules = append(form.Rules, flowcmp.RuleView{
			Field: flowcmp.Field{
				Source: fieldNames[i],
				Value:  fieldValue,
			},
			Condition: conditionTypes[i],
			Value:     conditionValues[i],
		})
	}

	form.RequireAll, err = strconv.ParseBool(req.RequireAll)
	if err != nil {
		return nil, errors.Wrap(err, "decodeFlowBuilderRequest: parsing requireAll")
	}

	return form, nil
}
