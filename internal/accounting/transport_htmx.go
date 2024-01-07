package accounting

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"

	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/itsnoproblem/prmry/internal/components"
	"github.com/itsnoproblem/prmry/internal/htmx"
	internalhttp "github.com/itsnoproblem/prmry/internal/http"
)

type Renderer interface {
	Render(w http.ResponseWriter, r *http.Request, cmp components.Component) error
	RenderError(w http.ResponseWriter, r *http.Request, err error)
}

func RouteHandler(svc Service, renderer Renderer) func(chi.Router) {
	getAccountEndpoint := internalhttp.NewHTMXEndpoint(
		makeGetAccountEndpoint(svc),
		decodeGetAccountRequest,
		formatAccountResponse,
		auth.Required,
	)

	updateProfileEndpoint := internalhttp.NewHTMXEndpoint(
		makeUpdateProfileEndpoint(svc),
		decodeUpdateProfileRequest,
		formatUpdateProfileResponse,
		auth.Required,
	)

	createAPIKeyEndpoint := internalhttp.NewHTMXEndpoint(
		makeCreateAPIKeyEndpoint(svc),
		decodeEmptyRequest,
		formatAPIKeyResponse,
		auth.Required,
	)

	updateAPIKeyEndpoint := internalhttp.NewHTMXEndpoint(
		makeUpdateAPIKeyEndpoint(svc),
		decodeUpdateAPIKeyRequest,
		formatUpdateAPIKeyResponse,
		auth.Required,
	)

	deleteAPIKeyEndpoint := internalhttp.NewHTMXEndpoint(
		makeDeleteAPIKeyEndpoint(svc),
		decodeDeleteAPIKeyRequest,
		formatDeleteAPIKeyResponse,
		auth.Required,
	)

	return func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Get("/account", htmx.MakeHandler(getAccountEndpoint, renderer))
			r.Put("/account/profile", htmx.MakeHandler(updateProfileEndpoint, renderer))
			r.Post("/account/api-keys", htmx.MakeHandler(createAPIKeyEndpoint, renderer))
			r.Put("/account/api-keys/{keyID}", htmx.MakeHandler(updateAPIKeyEndpoint, renderer))
			r.Delete("/account/api-keys/{keyID}", htmx.MakeHandler(deleteAPIKeyEndpoint, renderer))
		})
	}
}

func decodeEmptyRequest(_ context.Context, request *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeGetAccountRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	usr := auth.UserFromContext(ctx)
	if usr == nil {
		return nil, fmt.Errorf("missing user")
	}

	return getAccountRequest{
		UserID: usr.ID,
	}, nil
}

type updateProfileRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func decodeUpdateProfileRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	usr := auth.UserFromContext(ctx)
	if usr == nil {
		return nil, fmt.Errorf("missing user")
	}

	if request.Body == nil {
		return nil, fmt.Errorf("decodeUpdateProfileRequest: missing request body")
	}

	var req updateProfileRequest
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, errors.Wrap(err, "decodeUpdateProfileRequest")
	}

	if err = json.Unmarshal(body, &req); err != nil {
		return nil, errors.Wrap(err, "decodeUpdateProfileRequest")
	}

	return profileOptions{
		UserID: usr.ID,
		Name:   req.Name,
		Email:  usr.Email, // don't allow changing email for now
	}, nil
}

type updateAPIKeyRequest struct {
	KeyID string
	Name  string `json:"keyName"`
}

func decodeUpdateAPIKeyRequest(_ context.Context, request *http.Request) (interface{}, error) {
	if request.Body == nil {
		return nil, fmt.Errorf("decodeUpdateAPIKeyRequest: missing request body")
	}

	var req updateAPIKeyRequest
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, errors.Wrap(err, "decodeUpdateAPIKeyRequest")
	}

	if err = json.Unmarshal(body, &req); err != nil {
		return nil, errors.Wrap(err, "decodeUpdateAPIKeyRequest")
	}

	req.KeyID = chi.URLParam(request, "keyID")
	return updateAPIKeyOptions{
		KeyID: req.KeyID,
		Name:  req.Name,
	}, nil
}

func decodeDeleteAPIKeyRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	return deleteAPIKeyRequest{
		KeyID: chi.URLParam(request, "keyID"),
	}, nil
}
