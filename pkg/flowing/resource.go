package flowing

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"

	flowcmp "github.com/itsnoproblem/prmry/pkg/components/flow"
	"github.com/itsnoproblem/prmry/pkg/components/success"
	"github.com/itsnoproblem/prmry/pkg/flow"
)

type Service interface {
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

	r.Post("/new", rs.CreateFlow)
	r.Get("/new", rs.FlowBuilder)
	r.Post("/new/add-rule", rs.FlowBuilderAddRule)

	//
	r.Get("/", rs.ListFlows)

	// Get a flow by ID
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", rs.GetFlow)
	})

	return r
}

// Flow Builder

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

func (rs Resource) FlowBuilderAddRule(w http.ResponseWriter, r *http.Request) {
	cmp, err := readForm(r)
	if err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}

	cmp.Conditions = append(cmp.Conditions, flowcmp.ConditionView{})
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

	frag := flowcmp.FlowBuilder(cmp)
	page := flowcmp.FlowBuilderPage(cmp)

	if err = rs.renderer.RenderTemplComponent(w, r, page, frag); err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}

	notice := success.SuccessView{
		Message: "Saved flow",
	}
	success.Success(notice).Render(r.Context(), w)
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
	FieldNames      interface{} `json:"fieldName"`
	ConditionTypes  interface{} `json:"condition"`
	ConditionValues interface{} `json:"value"`
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

	if req.FieldNames != nil {
		v := reflect.ValueOf(req.FieldNames)
		switch v.Kind() {
		case reflect.String:
			fieldNames = append(fieldNames, req.FieldNames.(string))
			conditionTypes = append(conditionTypes, req.ConditionTypes.(string))
			conditionValues = append(conditionValues, req.ConditionValues.(string))
			break

		case reflect.Slice:
			fieldNames, err = stringSlice(req.FieldNames)
			if err != nil {
				return flowcmp.Detail{}, errors.Wrap(err, "readForm: fieldNames")
			}

			conditionTypes, err = stringSlice(req.ConditionTypes)
			if err != nil {
				return flowcmp.Detail{}, errors.Wrap(err, "readForm: conditionTypes")
			}

			conditionValues, err = stringSlice(req.ConditionValues)
			if err != nil {
				return flowcmp.Detail{}, errors.Wrap(err, "readForm: conditionValues")
			}

			break

		default:
			return flowcmp.Detail{}, fmt.Errorf("readForm: FieldNames unknown type %T", fieldNames)
		}
	}

	if len(fieldNames) != len(conditionTypes) || len(fieldNames) != len(conditionValues) {
		return flowcmp.Detail{}, fmt.Errorf("readForm: condition fields mismatch")
	}

	conditions := make([]flowcmp.ConditionView, 0)
	for i, _ := range fieldNames {
		conditions = append(conditions, flowcmp.ConditionView{
			Type:  conditionTypes[i],
			Field: fieldNames[i],
			Value: conditionValues[i],
		})
	}

	return flowcmp.Detail{
		ID:         req.ID,
		Name:       req.Name,
		Conditions: conditions,
	}, nil
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
