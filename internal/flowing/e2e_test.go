package flowing_test

import (
	"context"
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

func TestFlowing(t *testing.T) {
	renderer := htmx.NewRenderer()
	flowsRepo := inmem.NewFlowsRepo()
	flowSvc := flowing.NewService(flowsRepo)

	r := chi.NewRouter()
	r.Use(auth.Middleware(auth.Byte32{}, false))
	r.Use(render.SetContentType(render.ContentTypeHTML))
	r.Use(htmx.Middleware)
	r.Use(auth.TestUserMiddleware)

	r.Route("/", flowing.RouteHandler(flowSvc, renderer))

	server := httptest.NewServer(r)
	defer server.Close()
	client := server.Client()
	user := auth.TestModeUser()

	if err := createSomeFlows(context.Background(), user.ID, flowsRepo); err != nil {
		t.Fatalf(err.Error())
	}

	testCases := []struct {
		Name           string
		Method         string
		Path           string
		FullPage       bool
		WantStatusCode int
	}{
		{
			Name:           "New Flow Form Page",
			Method:         http.MethodGet,
			Path:           "/flows/new",
			FullPage:       true,
			WantStatusCode: http.StatusOK,
		},
		{
			Name:           "New Flow Form Fragment",
			Method:         http.MethodGet,
			Path:           "/flows/new",
			FullPage:       false,
			WantStatusCode: http.StatusOK,
		},
	}

	for _, tt := range testCases {

		t.Run(tt.Name, func(t *testing.T) {
			u, err := url.Parse(server.URL + tt.Path)
			if err != nil {
				t.Fatalf("parse url: %s", err)
			}

			req := &http.Request{
				Method: http.MethodGet,
				URL:    u,
			}

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

			approvals.VerifyString(t, string(data[:]))
		})
	}

}

func createSomeFlows(ctx context.Context, userID string, repo flowing.FlowsRepo) error {
	flows := []flow.Flow{
		flow.Flow{
			ID:     "123",
			UserID: userID,
			Name:   "Test Flow A",
		},
		flow.Flow{
			ID:     "567",
			UserID: userID,
			Name:   "Test Flow A",
		},
		flow.Flow{
			ID:     "998",
			UserID: userID,
			Name:   "Test Flow A",
		},
	}

	for _, f := range flows {
		if err := repo.InsertFlow(ctx, f); err != nil {
			return errors.Wrap(err, "createSomeFlows")
		}
	}

	return nil
}
