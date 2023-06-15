package flowing

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/itsnoproblem/prmry/pkg/components/success"
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
	GetFlowsForUser(ctx context.Context, userID string) ([]flow.Flow, error)
	CreateFlow(ctx context.Context, flw flow.Flow) (ID string, err error)
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

	r.Post("/new", rs.CreateFlow)
	r.Get("/new", rs.FlowBuilder)
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

// FlowBuilder GET /flows/new
func (rs Resource) FlowBuilder(w http.ResponseWriter, r *http.Request) {
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

// FlowBuilderAddRule POST /flows/new/rules
func (rs Resource) FlowBuilderAddRule(w http.ResponseWriter, r *http.Request) {
	cmp, err := readForm(r)
	if err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}

	cmp.Rules = append(cmp.Rules, flowcmp.RuleView{})
	rs.FlowBuilder(w, r.WithContext(context.WithValue(r.Context(), "view", cmp)))
}

// FlowBuilderRemoveRule DELETE /flows/new/rules/{index}
func (rs Resource) FlowBuilderRemoveRule(w http.ResponseWriter, r *http.Request) {
	index := chi.URLParam(r, "index")
	if index == "" {
		rs.renderer.RenderError(w, r, fmt.Errorf("missing required parameter: index"))
		return
	}

	idx, err := strconv.Atoi(index)
	if err != nil {
		rs.renderer.RenderError(w, r, errors.Wrap(err, "FlowBuilderRemoveRule"))
	}

	cmp, err := readForm(r)
	if err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}

	conditions := make([]flowcmp.RuleView, 0)
	conditions = append(conditions, cmp.Rules[:idx]...)
	cmp.Rules = append(conditions, cmp.Rules[idx+1:]...)

	rs.FlowBuilder(w, r.WithContext(context.WithValue(r.Context(), "view", cmp)))
}

// FlowBuilderUpdatePrompt PUT /flows/new/prompt
func (rs Resource) FlowBuilderUpdatePrompt(w http.ResponseWriter, r *http.Request) {
	cmp, err := readForm(r)
	if err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}

	re, err := regexp.Compile("%s")
	if err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}

	matches := re.FindAllString(cmp.Prompt, -1)
	if len(matches) > len(cmp.PromptArgs) {
		for i := len(cmp.PromptArgs); i < len(matches); i++ {
			cmp.PromptArgs = append(cmp.PromptArgs, "")
		}
	}

	if len(matches) < len(cmp.PromptArgs) {
		cmp.PromptArgs = cmp.PromptArgs[:len(matches)]
	}

	rs.FlowBuilder(w, r.WithContext(context.WithValue(r.Context(), "view", cmp)))
}

// CreateFlow - POST /flows
func (rs Resource) CreateFlow(w http.ResponseWriter, r *http.Request) {
	cmp, err := readForm(r)
	if err != nil {
		rs.renderer.RenderError(w, r, err)
		rs.FlowBuilder(w, r)
		return
	}

	if err = cmp.Lock(r); err != nil {
		rs.renderer.RenderError(w, r, err)
		rs.FlowBuilder(w, r.WithContext(context.WithValue(r.Context(), "view", cmp)))
		return
	}

	flw := cmp.ToFlow()
	ok := true
	if cmp.ID, err = rs.service.CreateFlow(r.Context(), flw); err != nil {
		rs.renderer.RenderError(w, r, err)
		ok = false
	}

	frag := flowcmp.FlowBuilder(cmp)
	page := flowcmp.FlowBuilderPage(cmp)

	if ok {
		notice := success.SuccessView{
			Message: "Saved " + flw.Name,
		}
		if err = success.Success(notice).Render(r.Context(), w); err != nil {
			rs.renderer.RenderError(w, r, err)
		}
	}

	if err = rs.renderer.RenderTemplComponent(w, r, page, frag); err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}
}

// ListFlows - GET /flows
func (rs Resource) ListFlows(w http.ResponseWriter, r *http.Request) {

	cmp := flowcmp.FlowsListView{
		Flows: []flowcmp.FlowSummary{
			{
				ID:   "1234-fake",
				Name: "Mr. Clean Campaign",
			},
			{
				ID:   "7846-two",
				Name: "Sarcastic Robot",
			},
		},
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
	ConditionTypes  interface{} `json:"condition"`
	ConditionValues interface{} `json:"value"`
	Prompt          string      `json:"prompt"`
	PromptArgs      interface{} `json:"promptArgs"`
}

func readForm(r *http.Request) (flowcmp.Detail, error) {
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
	conditionTypes := make([]string, 0)
	conditionValues := make([]string, 0)
	responseArgs := make([]string, 0)

	fieldNames, err = stringOrSlice(req.FieldNames)
	if err != nil {
		return flowcmp.Detail{}, errors.Wrap(err, "readForm")
	}

	conditionTypes, err = stringOrSlice(req.ConditionTypes)
	if err != nil {
		return flowcmp.Detail{}, errors.Wrap(err, "readForm")
	}

	conditionValues, err = stringOrSlice(req.ConditionValues)
	if err != nil {
		return flowcmp.Detail{}, errors.Wrap(err, "readForm")
	}

	responseArgs, err = stringOrSlice(req.PromptArgs)
	if err != nil {
		return flowcmp.Detail{}, errors.Wrap(err, "readForm")
	}

	if len(fieldNames) != len(conditionTypes) || len(fieldNames) != len(conditionValues) {
		return flowcmp.Detail{}, fmt.Errorf("readForm: condition fields mismatch")
	}

	conditions := make([]flowcmp.RuleView, 0)
	for i, _ := range fieldNames {
		conditions = append(conditions, flowcmp.RuleView{
			Condition: conditionTypes[i],
			Field:     fieldNames[i],
			Value:     conditionValues[i],
		})
	}

	requireAll, err := strconv.ParseBool(req.RequireAll)
	if err != nil {
		return flowcmp.Detail{}, fmt.Errorf("readForm: parsing requireAll")
	}

	return flowcmp.Detail{
		ID:                  req.ID,
		Name:                req.Name,
		Rules:               conditions,
		RequireAll:          requireAll,
		Prompt:              req.Prompt,
		PromptArgs:          responseArgs,
		SupportedFields:     flow.SupportedFields(),
		SupportedConditions: flow.SupportedConditions(),
	}, nil
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
