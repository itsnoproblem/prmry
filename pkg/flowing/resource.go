package flowing

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/itsnoproblem/prmry/pkg/auth"
	"github.com/itsnoproblem/prmry/pkg/htmx"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"strconv"

	flowcmp "github.com/itsnoproblem/prmry/pkg/components/flow"
	"github.com/itsnoproblem/prmry/pkg/flow"
)

type Service interface {
	CreateFlow(ctx context.Context, flw flow.Flow) (ID string, err error)
	UpdateFlow(ctx context.Context, flw flow.Flow) error
	GetFlow(ctx context.Context, flowID string) (flow.Flow, error)
	GetFlowsForUser(ctx context.Context, userID string) ([]flow.Flow, error)
}

type Renderer interface {
	RenderError(w http.ResponseWriter, r *http.Request, err error)
	RenderTemplComponent(w http.ResponseWriter, r *http.Request, fullPage, fragment templ.Component) error
	Unauthorized(w http.ResponseWriter, r *http.Request)
}

type Resource struct {
	renderer Renderer
	service  Service
}

func NewResource(renderer Renderer, svc Service) *Resource {
	return &Resource{
		renderer: renderer,
		service:  svc,
	}
}

func (rs Resource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/{flowID}/edit", rs.EditFlowForm)
	r.Put("/{flowID}", rs.SaveFlow)

	r.Post("/", rs.SaveFlow)
	r.Get("/new", rs.NewFlowForm)
	r.Put("/new/prompt", rs.FlowBuilderUpdatePrompt)
	r.Post("/new/rules", rs.FlowBuilderAddRule)
	r.Delete("/new/rules/{index}", rs.FlowBuilderRemoveRule)

	//
	r.Get("/", rs.ListFlows)

	// Get a flow by ID
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", rs.GetFlow)
	})

	return r
}

// NewFlowForm GET /flows/new
func (rs Resource) NewFlowForm(w http.ResponseWriter, r *http.Request) {
	cmp := flowcmp.Detail{
		SupportedFields:     flow.SupportedFields(),
		SupportedConditions: flow.SupportedConditions(),
	}
	if existing := r.Context().Value("view"); existing != nil {
		cmp = existing.(flowcmp.Detail)
	}

	if err := cmp.Lock(r); err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}

	page := flowcmp.FlowBuilderPage(cmp)
	fragment := flowcmp.FlowBuilder(cmp)

	if err := rs.renderer.RenderTemplComponent(w, r, page, fragment); err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}
}

// EditFlowForm GET /flows/{flowID}/edit
func (rs Resource) EditFlowForm(w http.ResponseWriter, r *http.Request) {
	var cmp flowcmp.Detail
	ctx := r.Context()

	flowID := chi.URLParam(r, "flowID")
	if flowID == "" {
		rs.renderer.RenderError(w, r, fmt.Errorf("missing flowID"))
		return
	}

	if existing := r.Context().Value("view"); existing != nil {
		cmp = existing.(flowcmp.Detail)
	} else {
		flw, err := rs.service.GetFlow(ctx, flowID)
		if err != nil {
			rs.renderer.RenderError(w, r, err)
		}

		cmp = flowcmp.NewDetail(flw)
	}

	if err := cmp.Lock(r); err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}

	user := auth.UserFromContext(ctx)
	if user == nil {
		rs.renderer.RenderError(w, r, fmt.Errorf("missing user"))
	}

	var err error
	cmp.AvailableFlowsByID, err = rs.getAvailableFlows(ctx)
	if err != nil {
		rs.renderer.RenderError(w, r, err)
	}

	page := flowcmp.FlowBuilderPage(cmp)
	fragment := flowcmp.FlowBuilder(cmp)

	if err := rs.renderer.RenderTemplComponent(w, r, page, fragment); err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}
}

// FlowBuilderAddRule POST /flows/new/rules
func (rs Resource) FlowBuilderAddRule(w http.ResponseWriter, r *http.Request) {
	cmp, err := rs.readForm(r)
	if err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}

	cmp.Rules = append(cmp.Rules, flowcmp.RuleView{})
	rs.NewFlowForm(w, r.WithContext(context.WithValue(r.Context(), "view", cmp)))
}

// FlowBuilderRemoveRule DELETE /flows/new/rules/{index}
func (rs Resource) FlowBuilderRemoveRule(w http.ResponseWriter, r *http.Request) {
	index := chi.URLParam(r, "index")
	if index == "" {
		rs.renderer.RenderError(w, r, fmt.Errorf("missing required parameter: index"))
	}

	idx, err := strconv.Atoi(index)
	if err != nil {
		rs.renderer.RenderError(w, r, errors.Wrap(err, "FlowBuilderRemoveRule"))
	}

	cmp, err := rs.readForm(r)
	if err != nil {
		rs.renderer.RenderError(w, r, err)
	}

	conditions := make([]flowcmp.RuleView, 0)
	conditions = append(conditions, cmp.Rules[:idx]...)
	cmp.Rules = append(conditions, cmp.Rules[idx+1:]...)

	rs.NewFlowForm(w, r.WithContext(context.WithValue(r.Context(), "view", cmp)))
}

// FlowBuilderUpdatePrompt PUT /flows/new/prompt
func (rs Resource) FlowBuilderUpdatePrompt(w http.ResponseWriter, r *http.Request) {
	cmp, err := rs.readForm(r)
	if err != nil {
		rs.renderer.RenderError(w, r, err)
	}

	re, err := regexp.Compile("%s")
	if err != nil {
		rs.renderer.RenderError(w, r, err)
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

	rs.NewFlowForm(w, r.WithContext(context.WithValue(r.Context(), "view", cmp)))
}

// SaveFlow - POST /flows | PUT /flows/{flowID}
func (rs Resource) SaveFlow(w http.ResponseWriter, r *http.Request) {
	cmp, err := rs.readForm(r)
	if err != nil {
		rs.renderer.RenderError(w, r, err)
		rs.NewFlowForm(w, r)
		return
	}

	if err = cmp.Lock(r); err != nil {
		rs.renderer.RenderError(w, r, err)
		rs.NewFlowForm(w, r.WithContext(context.WithValue(r.Context(), "view", cmp)))
		return
	}

	flw := cmp.ToFlow()
	ok := true

	if cmp.ID == "" {
		if cmp.ID, err = rs.service.CreateFlow(r.Context(), flw); err != nil {
			rs.renderer.RenderError(w, r, err)
			ok = false
		}
	} else {
		if err = rs.service.UpdateFlow(r.Context(), flw); err != nil {
			rs.renderer.RenderError(w, r, err)
			ok = false
		}
	}

	if !ok {
		frag := flowcmp.FlowBuilder(cmp)
		page := flowcmp.FlowBuilderPage(cmp)

		if err = rs.renderer.RenderTemplComponent(w, r, page, frag); err != nil {
			rs.renderer.RenderError(w, r, err)
			return
		}
	}

	htmx.Redirect(w, "/flows")
	return
}

// ListFlows - GET /flows
func (rs Resource) ListFlows(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		rs.renderer.RenderError(w, r, fmt.Errorf("user is missing"))
	}

	flows, err := rs.service.GetFlowsForUser(r.Context(), user.ID)
	if err != nil {
		rs.renderer.RenderError(w, r, err)
	}

	summaries := make([]flowcmp.FlowSummary, 0)
	for _, flow := range flows {
		label := "rule"
		if len(flow.Rules) > 1 {
			label = "rules"
		}

		summaries = append(summaries, flowcmp.FlowSummary{
			ID:          flow.ID,
			Name:        flow.Name,
			RuleCount:   fmt.Sprintf("%d %s", len(flow.Rules), label),
			LastChanged: flow.UpdatedAt.Format("Jan 02, 2006 15:04"),
		})
	}

	cmp := flowcmp.FlowsListView{
		Flows: summaries,
	}

	if err := cmp.Lock(r); err != nil {
		rs.renderer.Unauthorized(w, r)
		return
	}

	page := flowcmp.FlowsListPage(cmp)
	fragment := flowcmp.FlowsList(cmp)

	if err := rs.renderer.RenderTemplComponent(w, r, page, fragment); err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}
}

// GetFlow - GET /flows/{id}
func (rs Resource) GetFlow(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	if id == "" {
		rs.renderer.RenderError(w, r, fmt.Errorf("missing required parameter 'id'"))
		return
	}

	cmp := flowcmp.Detail{
		ID:   "fake-1234",
		Name: "Test Fake Flow",
	}

	if err := cmp.Lock(r); err != nil {
		rs.renderer.Unauthorized(w, r)
		return
	}

	page := flowcmp.FlowBuilderPage(cmp)
	fragment := flowcmp.FlowBuilder(cmp)

	if err := rs.renderer.RenderTemplComponent(w, r, page, fragment); err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}
}

// --- private ---

type flowBuilderForm struct {
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

func (rs Resource) readForm(r *http.Request) (flowcmp.Detail, error) {
	if r == nil {
		return flowcmp.Detail{}, fmt.Errorf("readForm: request was null")
	}

	var req flowBuilderForm
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return flowcmp.Detail{}, errors.Wrap(err, "readForm")
	}

	if err = json.Unmarshal(body, &req); err != nil {
		return flowcmp.Detail{}, errors.Wrap(err, "readForm")
	}

	fieldNames := make([]string, 0)
	selectedFlows := make([]string, 0)
	conditionTypes := make([]string, 0)
	conditionValues := make([]string, 0)
	parsedPromptArgs := make([]string, 0)
	promptArgFlows := make([]string, 0)

	fieldNames, err = stringOrSlice(req.FieldNames)
	if err != nil {
		return flowcmp.Detail{}, errors.Wrap(err, "readForm: fieldNames")
	}

	selectedFlows, err = stringOrSlice(req.SelectedFlows)
	if err != nil {
		return flowcmp.Detail{}, errors.Wrap(err, "readForm: selectedFlows")
	}

	conditionTypes, err = stringOrSlice(req.ConditionTypes)
	if err != nil {
		return flowcmp.Detail{}, errors.Wrap(err, "readForm: conditionTypes")
	}

	conditionValues, err = stringOrSlice(req.ConditionValues)
	if err != nil {
		return flowcmp.Detail{}, errors.Wrap(err, "readForm: conditionValues")
	}

	parsedPromptArgs, err = stringOrSlice(req.PromptArgs)
	if err != nil {
		return flowcmp.Detail{}, errors.Wrap(err, "readForm: parsedPromptArgs")
	}

	if req.PromptArgFlows != nil {
		promptArgFlows, err = stringOrSlice(req.PromptArgFlows)
		if err != nil {
			return flowcmp.Detail{}, errors.Wrap(err, "readForm: promptArgFlows")
		}
	}

	if len(fieldNames) != len(conditionTypes) || len(fieldNames) != len(conditionValues) {
		return flowcmp.Detail{}, fmt.Errorf("readForm: condition fields mismatch")
	}

	availableFlowsByID, err := rs.getAvailableFlows(r.Context())
	if err != nil {
		return flowcmp.Detail{}, errors.Wrap(err, "readForm")
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
		AvailableFlowsByID:  availableFlowsByID,
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
		return flowcmp.Detail{}, fmt.Errorf("readForm: parsing requireAll")
	}

	return form, nil
}

func (rs Resource) getAvailableFlows(ctx context.Context) (map[string]string, error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("getAvailableFlows: missing user")
	}

	flows, err := rs.service.GetFlowsForUser(ctx, user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "getAvailableFlows")
	}

	availableFlowsByID := make(map[string]string)
	for _, flw := range flows {
		availableFlowsByID[flw.ID] = flw.Name
	}

	return availableFlowsByID, nil
}

func stringOrSlice(str interface{}) ([]string, error) {
	var err error
	result := make([]string, 0)
	v := reflect.ValueOf(str)

	switch v.Kind() {
	case reflect.String:
		result = append(result, str.(string))
		break

	case reflect.Slice:
		result, err = stringSlice(str)
		if err != nil {
			return nil, errors.Wrap(err, "stringOrSlice")
		}
		break

	default:
		return nil, nil
	}

	return result, nil
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
