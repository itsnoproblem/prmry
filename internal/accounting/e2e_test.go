package accounting_test

import (
	"context"
	"github.com/approvals/go-approval-tests/reporters"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/approvals/go-approval-tests"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/itsnoproblem/prmry/internal/accounting"
	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/itsnoproblem/prmry/internal/htmx"
)

func init() {
	approvals.UseFolder("testdata")
}

func TestAccounting(t *testing.T) {
	reporter := approvals.UseReporter(reporters.NewFileMergeReporter())
	defer reporter.Close()

	renderer := htmx.NewRenderer()
	service := accounting.NewService(&userRepoMock{})

	r := chi.NewRouter()
	r.Use(auth.Middleware(auth.Byte32{}, false))
	r.Use(render.SetContentType(render.ContentTypeHTML))
	r.Use(htmx.Middleware)
	r.Use(auth.TestUserMiddleware)

	r.Route("/", accounting.RouteHandler(service, renderer))

	server := httptest.NewServer(r)
	defer server.Close()
	client := server.Client()

	testCases := []struct {
		Name           string
		Method         string
		Path           string
		FullPage       bool
		WantStatusCode int
	}{
		{
			Name:           "My Account Page",
			Method:         http.MethodGet,
			Path:           "/account",
			FullPage:       true,
			WantStatusCode: http.StatusOK,
		},
		{
			Name:           "My Account Fragment",
			Method:         http.MethodGet,
			Path:           "/account",
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
				t.Errorf("expected response code %d but received %d", tt.WantStatusCode, resp.StatusCode)
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

type userRepoMock struct{} // accounting.UserRepo
func (m *userRepoMock) FindUserByID(ctx context.Context, id string) (usr auth.User, exists bool, err error) {
	return fakeUser(), true, nil
}
func (m *userRepoMock) UpdateAccountProfile(ctx context.Context, userID, name, email string) error {
	return nil
}
func (m *userRepoMock) FindAPIKeysForUser(ctx context.Context, userID string) ([]auth.APIKey, error) {
	return []auth.APIKey{
		{
			Key:  "123",
			Name: "Test Key",
		},
	}, nil
}
func (m *userRepoMock) UpdateAPIKeyName(ctx context.Context, userID, keyID, name string) error {
	return nil
}
func (m *userRepoMock) InsertAPIKey(ctx context.Context, userID string, key auth.APIKey) error {
	return nil
}
func (m *userRepoMock) DeleteAPIKey(ctx context.Context, userID, keyID string) error {
	return nil
}

func fakeUser() auth.User {
	return auth.User{
		ID:       "123",
		Email:    "example@prmry.io",
		Name:     "Test User",
		Nickname: "testy",
	}
}
