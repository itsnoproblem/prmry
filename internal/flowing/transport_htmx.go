package flowing

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"

	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/itsnoproblem/prmry/internal/components"
	flowcmp "github.com/itsnoproblem/prmry/internal/components/flow"
	"github.com/itsnoproblem/prmry/internal/flow"
	"github.com/itsnoproblem/prmry/internal/htmx"
	internalhttp "github.com/itsnoproblem/prmry/internal/http"
)

type Renderer interface {
	Render(w http.ResponseWriter, r *http.Request, cmp components.Component) error
	RenderError(w http.ResponseWriter, r *http.Request, err error)
}

func RouteHandler(svc Service, renderer Renderer) func(chi.Router) {
	listFlowsEndpoint := internalhttp.NewHTMXEndpoint(
		makeListFlowsEndpoint(svc),
		decodeEmptyRequest,
		formatFlowSummaries,
		auth.Required,
	)

	getFlowBuilderEndpoint := internalhttp.NewHTMXEndpoint(
		makeFlowBuilderEndpoint(svc),
		decodeEmptyRequest,
		formatFlowBuilderResponse,
		auth.Required,
	)

	updateFlowBuilderEndpoint := internalhttp.NewHTMXEndpoint(
		makeFlowBuilderEndpoint(svc),
		decodeFlowBuilderRequest,
		formatFlowBuilderResponse,
		auth.Required,
	)

	editFlowEndpoint := internalhttp.NewHTMXEndpoint(
		makeEditFlowFormEndpoint(svc),
		decodeEditFormRequest,
		formatFlowBuilderResponse,
		auth.Required,
	)

	flowBuilderAddRuleEndpoint := internalhttp.NewHTMXEndpoint(
		makeFlowBuilderAddRuleEndpoint(svc),
		decodeFlowBuilderRequest,
		formatFlowBuilderResponse,
		auth.Required,
	)

	flowBuilderDeleteRuleEndpoint := internalhttp.NewHTMXEndpoint(
		makeFlowBuilderRemoveRuleEndpoint(svc),
		decodeFlowBuilderDeleteRuleRequest,
		formatFlowBuilderResponse,
		auth.Required,
	)

	flowBuilderAddInputEndpoint := internalhttp.NewHTMXEndpoint(
		makeFlowBuilderAddInputEndpoint(svc),
		decodeFlowBuilderRequest,
		formatFlowBuilderResponse,
		auth.Required,
	)

	flowBuilderDeleteInputEndpoint := internalhttp.NewHTMXEndpoint(
		makeFlowBuilderRemoveInputEndpoint(svc),
		decodeFlowBuilderDeleteInputRequest,
		formatFlowBuilderResponse,
		auth.Required,
	)

	flowBuilderUpdatePromptEndpoint := internalhttp.NewHTMXEndpoint(
		makeFlowBuilderUpdatePromptEndpoint(svc),
		decodeFlowBuilderRequest,
		formatFlowBuilderPromptResponse,
		auth.Required,
	)

	saveFlowEndpoint := internalhttp.NewHTMXEndpoint(
		makeSaveFlowEndpoint(svc),
		decodeFlowBuilderRequest,
		formatRedirectResponse,
		auth.Required,
	)

	deleteFlowEndpoint := internalhttp.NewHTMXEndpoint(
		makeDeleteFlowEndpoint(svc),
		decodeDeleteFlowRequest,
		formatRedirectResponse,
		auth.Required,
	)

	return func(r chi.Router) {
		r.Group(func(r chi.Router) {

			r.Get("/flows", htmx.MakeHandler(listFlowsEndpoint, renderer))
			r.Get("/flows/{flowID}/edit", htmx.MakeHandler(editFlowEndpoint, renderer))
			r.Post("/flows", htmx.MakeHandler(saveFlowEndpoint, renderer))
			r.Delete("/flows/{flowID}", htmx.MakeHandler(deleteFlowEndpoint, renderer))

			r.Get("/flow-builder", htmx.MakeHandler(getFlowBuilderEndpoint, renderer))
			r.Put("/flow-builder", htmx.MakeHandler(updateFlowBuilderEndpoint, renderer))

			r.Put("/flow-builder/prompt", htmx.MakeHandler(flowBuilderUpdatePromptEndpoint, renderer))
			r.Post("/flow-builder/rules", htmx.MakeHandler(flowBuilderAddRuleEndpoint, renderer))
			r.Delete("/flow-builder/rules/{index}", htmx.MakeHandler(flowBuilderDeleteRuleEndpoint, renderer))
			r.Post("/flow-builder/inputs", htmx.MakeHandler(flowBuilderAddInputEndpoint, renderer))
			r.Delete("/flow-builder/inputs/{index}", htmx.MakeHandler(flowBuilderDeleteInputEndpoint, renderer))
		})
	}
}

func decodeEmptyRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeDeleteFlowRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	return deleteFlowRequest{
		FlowID:      chi.URLParam(request, "flowID"),
		SelectedTab: selectedTabFromURL(request),
	}, nil
}

func decodeEditFormRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	return editFlowRequest{
		FlowID:      chi.URLParam(request, "flowID"),
		SelectedTab: selectedTabFromURL(request),
	}, nil
}

func decodeFlowBuilderDeleteRuleRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	index := chi.URLParam(request, "index")
	idx, err := strconv.Atoi(index)
	if err != nil {
		return nil, errors.Wrap(err, "decodeFlowBuilderDeleteRuleRequest")
	}

	form, err := decodeFlowBuilderRequest(ctx, request)
	if err != nil {
		return nil, errors.Wrap(err, "decodeFlowBuilderDeleteRuleRequest")
	}

	formRequest, ok := form.(flowcmp.Detail)
	if !ok {
		return nil, fmt.Errorf("decodeFlowBuilderDeleteRuleRequest: failed to parse form")
	}

	return flowBuilderRemoveItemRequest{
		Index: idx,
		Form:  formRequest,
	}, nil
}

func decodeFlowBuilderDeleteInputRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	index := chi.URLParam(request, "index")
	idx, err := strconv.Atoi(index)
	if err != nil {
		return nil, errors.Wrap(err, "decodeFlowBuilderDeleteInputRequest")
	}

	form, err := decodeFlowBuilderRequest(ctx, request)
	if err != nil {
		return nil, errors.Wrap(err, "decodeFlowBuilderDeleteInputRequest")
	}

	formRequest, ok := form.(flowcmp.Detail)
	if !ok {
		return nil, fmt.Errorf("decodeFlowBuilderDeleteInputRequest: failed to parse form")
	}

	return flowBuilderRemoveItemRequest{
		Index: idx,
		Form:  formRequest,
	}, nil
}

type FlowBuilderFormRequest struct {
	ID             string      `json:"id"`
	Name           string      `json:"name"`
	SelectedTab    string      `json:"selectedTab"`
	AvailableFlows interface{} `json:"availableFlows"`

	// Triggers
	RequireAll      string      `json:"requireAll"`
	FieldNames      interface{} `json:"fieldName"`
	SelectedFlows   interface{} `json:"selectedFlows"`
	ConditionTypes  interface{} `json:"condition"`
	ConditionValues interface{} `json:"value"`
	RuleFlows       interface{} `json:"ruleConditionFlows"`
	RuleInputParams interface{} `json:"ruleInputParams"`

	// Prompt
	Prompt          string      `json:"prompt"`
	PromptArgs      interface{} `json:"promptArgs"`
	PromptArgFlows  interface{} `json:"promptArgFlows"`
	PromptArgInputs interface{} `json:"promptArgInputs"`

	// Inputs
	InputParams         interface{} `json:"inputParams"`
	InputParamsRequired interface{} `json:"inputParamsRequired"`
}

func decodeFlowBuilderRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	if r == nil {
		return flowcmp.Detail{}, fmt.Errorf("readForm: request was null")
	}

	var req FlowBuilderFormRequest
	if r.Body != nil {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, errors.Wrap(err, "decodeFlowBuilderRequest")
		}

		if err = json.Unmarshal(body, &req); err != nil {
			return nil, errors.Wrap(err, "decodeFlowBuilderRequest")
		}
	}

	promptArgs, err := makePromptArgs(req)
	if err != nil {
		return flowcmp.Detail{}, errors.Wrap(err, "makeInputParams")
	}

	inputParams, err := makeInputParams(req)
	if err != nil {
		return flowcmp.Detail{}, errors.Wrap(err, "makeInputParams")
	}

	selectedTab := selectedTabFromURL(r)
	if selectedTab == "" {
		selectedTab = req.SelectedTab
	}

	form := flowcmp.Detail{
		ID:                  req.ID,
		Name:                req.Name,
		Prompt:              req.Prompt,
		PromptArgs:          promptArgs,
		SupportedFields:     components.SortedMap(flow.SupportedFields()),
		SupportedConditions: components.SortedMap(flow.SupportedConditions()),
		InputParams:         inputParams,
		SelectedTab:         selectedTab,
	}

	if form.Rules, err = makeRules(req, inputParams); err != nil {
		return nil, errors.Wrap(err, "decodeFlowBuilderRequest")
	}

	form.RequireAll, err = strconv.ParseBool(req.RequireAll)
	if err != nil {
		form.RequireAll = false
	}

	return form, nil
}

func makeRules(req FlowBuilderFormRequest, inputParams []flowcmp.InputParam) ([]flowcmp.RuleView, error) {
	var err error
	flowIndex := 0
	paramIndex := 0
	conditionTypes := make([]string, 0)
	conditionValues := make([]string, 0)
	fieldNames := make([]string, 0)
	selectedFlows := make([]string, 0)

	fieldNames, err = htmx.StringOrSlice(req.FieldNames)
	if err != nil {
		return nil, errors.Wrap(err, "decodeFlowBuilderRequest: fieldNames")
	}

	selectedFlows, err = htmx.StringOrSlice(req.SelectedFlows)
	if err != nil {
		return nil, errors.Wrap(err, "decodeFlowBuilderRequest: selectedFlows")
	}

	ruleInputParams, err := htmx.StringOrSlice(req.RuleInputParams)
	if err != nil {
		return nil, errors.Wrap(err, "decodeFlowBuilderRequest: ruleInputParams")
	}

	conditionTypes, err = htmx.StringOrSlice(req.ConditionTypes)
	if err != nil {
		return nil, errors.Wrap(err, "decodeFlowBuilderRequest: conditionTypes")
	}

	conditionValues, err = htmx.StringOrSlice(req.ConditionValues)
	if err != nil {
		return nil, errors.Wrap(err, "decodeFlowBuilderRequest: conditionValues")
	}

	if len(fieldNames) != len(conditionTypes) || len(fieldNames) != len(conditionValues) {
		return nil, fmt.Errorf("decodeFlowBuilderRequest: condition fields mismatch")
	}

	rules := make([]flowcmp.RuleView, 0)
	for i, fieldSource := range fieldNames {
		fieldValue := ""
		if fieldSource == flow.FieldSourceFlow.String() && len(selectedFlows) > flowIndex {
			fieldValue = selectedFlows[flowIndex]
			flowIndex++
		}

		if fieldSource == flow.FieldSourceInputArg.String() && len(ruleInputParams) > paramIndex {
			fieldValue = ruleInputParams[paramIndex]
			paramIndex++
		}

		rules = append(rules, flowcmp.RuleView{
			Field: flowcmp.Field{
				Source: fieldNames[i],
				Value:  fieldValue,
			},
			Condition: conditionTypes[i],
			Value:     conditionValues[i],
		})
	}

	return rules, nil
}

func makeInputParams(req FlowBuilderFormRequest) ([]flowcmp.InputParam, error) {
	parsedInputParams, err := htmx.StringOrSlice(req.InputParams)
	if err != nil {
		return nil, errors.Wrap(err, "decodeFlowBuilderRequest: parsedInputParams")
	}

	parsedInputParamsRequired, err := htmx.StringOrSlice(req.InputParamsRequired)
	if err != nil {
		return nil, errors.Wrap(err, "decodeFlowBuilderRequest: parsedInputParamsRequired")
	}

	inputParams := make([]flowcmp.InputParam, 0)
	for i, param := range parsedInputParams {
		isRequired, err := strconv.ParseBool(parsedInputParamsRequired[i])
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse inputParams %d", i)
		}

		inputParam := flowcmp.InputParam{
			Type:     flow.ParamTypeString, // TODO(marty): un-hardcode this
			Key:      param,
			Required: isRequired,
		}

		inputParams = append(inputParams, inputParam)
	}

	return inputParams, nil
}

func makePromptArgs(req FlowBuilderFormRequest) ([]flowcmp.PromptArg, error) {
	parsedPromptArgs, err := htmx.StringOrSlice(req.PromptArgs)
	if err != nil {
		return nil, errors.Wrap(err, "makePromptArgs: parsedPromptArgs")
	}

	promptArgs := make([]flowcmp.PromptArg, 0)
	promptArgFlows := make([]string, 0)
	promptArgInputs := make([]string, 0)
	flowIndex := 0
	inputIndex := 0

	if req.PromptArgFlows != nil {
		promptArgFlows, err = htmx.StringOrSlice(req.PromptArgFlows)
		if err != nil {
			return nil, errors.Wrap(err, "decodeFlowBuilderRequest: promptArgFlows")
		}
	}

	if req.PromptArgInputs != nil {
		promptArgInputs, err = htmx.StringOrSlice(req.PromptArgInputs)
		if err != nil {
			return nil, errors.Wrap(err, "decodeFlowBuilderRequest: promptArgInputs")
		}
	}

	for _, arg := range parsedPromptArgs {
		pargs := flowcmp.PromptArg{
			Source: flow.SourceType(arg),
		}

		if flow.SourceType(arg) == flow.FieldSourceFlow && len(promptArgFlows) > flowIndex {
			pargs.Value = promptArgFlows[flowIndex]
			flowIndex++
		}

		if flow.SourceType(arg) == flow.FieldSourceInputArg && len(promptArgInputs) > inputIndex {
			pargs.Value = promptArgInputs[inputIndex]
			inputIndex++
		}

		promptArgs = append(promptArgs, pargs)
	}

	return promptArgs, nil
}

func selectedTabFromURL(r *http.Request) string {
	return r.URL.Query().Get("tab")
}

func stringSlice(input interface{}) ([]string, error) {
	interfaces, ok := input.([]interface{})
	if !ok {
		return nil, fmt.Errorf("stringSlice: failed to cast input to slice")
	}

	result := make([]string, len(interfaces))
	for i, val := range interfaces {
		strVal, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("stringSlice: failed to cast value to string")
		}

		result[i] = strVal
	}

	return result, nil
}
