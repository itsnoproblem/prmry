package htmx

import (
	"net/http"

	"github.com/google/martian/log"

	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/itsnoproblem/prmry/internal/components"
	"github.com/itsnoproblem/prmry/internal/components/redirect"
	internalhttp "github.com/itsnoproblem/prmry/internal/http"
)

type Renderer interface {
	Render(w http.ResponseWriter, r *http.Request, cmp components.Component) error
	RenderError(w http.ResponseWriter, r *http.Request, err error)
}

func MakeHandler(endpoint internalhttp.HTMXEndpoint, renderer Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		decoded, err := endpoint.DecodeRequest(ctx, r)
		if err != nil {
			log.Errorf("http.MakeHandler: %s", err)
			renderer.RenderError(w, r, err)
			return
		}

		user := auth.UserFromContext(ctx)
		if endpoint.RequiresAuth && user == nil {
			Redirect(w, "/")
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

		cmp, err := endpoint.EncodeResponse(ctx, res)
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
