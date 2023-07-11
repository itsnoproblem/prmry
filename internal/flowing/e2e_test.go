package flowing_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/approvals/go-approval-tests/reporters"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/approvals/go-approval-tests"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/pkg/errors"

	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/itsnoproblem/prmry/internal/flow"
	"github.com/itsnoproblem/prmry/internal/flowing"
	"github.com/itsnoproblem/prmry/internal/htmx"
	"github.com/itsnoproblem/prmry/internal/inmem"
)

func init() {
	approvals.UseFolder("testdata")
}

type testHelper struct {
	Renderer    flowing.Renderer
	FlowsRepo   flowing.FlowsRepo
	FlowService flowing.Service
	Router      *chi.Mux
}

func TestFlowing(t *testing.T) {
	reporter := approvals.UseReporter(reporters.NewGoLandReporter())
	defer reporter.Close()

	helper := newTestHelper()

	helper.Router.Route("/", flowing.RouteHandler(helper.FlowService, helper.Renderer))
	server := httptest.NewServer(helper.Router)
	defer server.Close()

	client := server.Client()
	user := auth.TestModeUser()

	if err := createSomeFlows(context.Background(), user.ID, helper.FlowsRepo); err != nil {
		t.Fatalf(err.Error())
	}

	testCases := []struct {
		Name           string
		Method         string
		Path           string
		FullPage       bool
		WantStatusCode int
		RequestPayload interface{}
	}{
		{
			Name:           "FlowBuilder Page",
			Method:         http.MethodGet,
			Path:           "/flows/new",
			FullPage:       true,
			WantStatusCode: http.StatusOK,
			RequestPayload: nil,
		},
		{
			Name:           "FlowBuilder Fragment",
			Method:         http.MethodGet,
			Path:           "/flows/new",
			FullPage:       false,
			WantStatusCode: http.StatusOK,
			RequestPayload: nil,
		},
		{
			Name:           "FlowBuilder Fragment AddRule",
			Method:         http.MethodPost,
			Path:           "/flows/new/rules",
			FullPage:       false,
			WantStatusCode: http.StatusOK,
			RequestPayload: requestPayloadNewRule{
				Id:         "",
				Name:       "",
				Prompt:     "",
				RequireAll: "false",
			},
		},
		{
			Name:           "FlowBuilder Fragment DeleteRule",
			Method:         http.MethodDelete,
			Path:           "/flows/new/rules/0",
			FullPage:       false,
			WantStatusCode: http.StatusOK,
			RequestPayload: requestPayloadNewRule{
				Id:         "",
				Name:       "",
				Prompt:     "",
				RequireAll: "false",
			},
		},
	}

	for _, tt := range testCases {

		t.Run(tt.Name, func(t *testing.T) {
			var reqPayload json.RawMessage

			u, err := url.Parse(server.URL + tt.Path)
			if err != nil {
				t.Fatalf("parse url: %s", err)
			}

			if tt.RequestPayload != nil {
				reqPayload, err = json.Marshal(tt.RequestPayload)
				if err != nil {
					t.Fatalf("encoding request payload: %s", err.Error())
				}
			}
			req, err := http.NewRequest(tt.Method, u.String(), bytes.NewReader(reqPayload))

			if !tt.FullPage {
				if req.Header == nil {
					req.Header = make(http.Header)
				}
				req.Header.Add(htmx.HeaderHXRequest, "true")
			}

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Get error: %v", err)
			}

			if resp.StatusCode != tt.WantStatusCode {
				t.Fatalf("expected response code %d but received %d", tt.WantStatusCode, resp.StatusCode)
			}

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("io.ReadAll error: %v", err)
			}

			extensionOpt := approvals.Options().WithExtension("html")
			approvals.VerifyString(t, string(data[:]), extensionOpt)
		})
	}

}

func newTestHelper() testHelper {
	r := chi.NewRouter()
	r.Use(auth.Middleware(auth.Byte32{}, false))
	r.Use(render.SetContentType(render.ContentTypeHTML))
	r.Use(htmx.Middleware)
	r.Use(auth.TestUserMiddleware)

	flowsRepo := inmem.NewFlowsRepo()

	return testHelper{
		Renderer:    htmx.NewRenderer(),
		FlowsRepo:   flowsRepo,
		FlowService: flowing.NewService(flowsRepo),
		Router:      r,
	}
}

type requestPayloadNewRule struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	RequireAll string `json:"requireAll"`
	Prompt     string `json:"prompt"`
}

func createSomeFlows(ctx context.Context, userID string, repo flowing.FlowsRepo) error {
	fakeRules := []flow.Rule{
		flow.Rule{
			Field: flow.Field{
				Source: "fakeSource",
				Value:  "fakeValue",
			},
			Condition: "equals",
			Value:     "test",
		},
	}

	flows := []flow.Flow{
		flow.Flow{
			ID:     "123",
			UserID: userID,
			Name:   "Test Flow A",
			Rules:  fakeRules,
		},
		flow.Flow{
			ID:     "567",
			UserID: userID,
			Name:   "Test Flow A",
			Rules:  fakeRules,
		},
		flow.Flow{
			ID:     "998",
			UserID: userID,
			Name:   "Test Flow A",
			Rules:  fakeRules,
		},
	}

	for _, f := range flows {
		if err := repo.InsertFlow(ctx, f); err != nil {
			return errors.Wrap(err, "createSomeFlows")
		}
	}

	return nil
}
