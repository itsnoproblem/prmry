package htmx

import (
	"fmt"
	"github.com/google/martian/log"
	"net/http"

	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/itsnoproblem/prmry/internal/components"
	"github.com/itsnoproblem/prmry/internal/components/redirect"
)

type Renderer interface {
	Render(w http.ResponseWriter, r *http.Request, cmp components.Component) error
	RenderError(w http.ResponseWriter, r *http.Request, err error)
}

func MakeHandler(endpoint Endpoint, renderer Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		decoded, err := endpoint.DecodeRequest(r)
		if err != nil {
			log.Errorf("http.MakeHandler: %s", err)
			renderer.RenderError(w, r, err)
			return
		}

		user := auth.UserFromContext(ctx)
		if endpoint.RequiresAuth && user == nil {
			renderer.RenderError(w, r, fmt.Errorf("you must be logged in to do that"))
			return
		}

		res, err := endpoint.HandleRequest(ctx, decoded)
		if err != nil {
			log.Errorf("http.MakeHandler: %s", err)
			renderer.RenderError(w, r, err)
			return
		}

		if redirect, isRedirect := res.(redirect.View); isRedirect {
			if IsHXRequest(r.Context()) {
				Redirect(w, redirect.Location)
				return
			}

			http.Redirect(w, r, redirect.Location, redirect.Status)
		}

		cmp, err := endpoint.EncodeResponse(res)
		if err != nil {
			log.Errorf("http.MakeHandler: %s", err)
			renderer.RenderError(w, r, err)
			return
		}

		if err = renderer.Render(w, r, cmp); err != nil {
			log.Errorf("http.MakeHandler: %s", err)
			renderer.RenderError(w, r, err)
		}
	}
}
