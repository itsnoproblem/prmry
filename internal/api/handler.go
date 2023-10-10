package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/martian/log"

	internalhttp "github.com/itsnoproblem/prmry/internal/http"
)

type Renderer interface {
	Render(w http.ResponseWriter, r *http.Request, data json.RawMessage) error
	RenderError(w http.ResponseWriter, r *http.Request, err error)
}

func MakeHandler(endpoint internalhttp.JSONEndpoint, renderer Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		decoded, err := endpoint.DecodeRequest(ctx, r)
		if err != nil {
			log.Errorf("http.MakeHandler: %s", err)
			renderer.RenderError(w, r, err)
			return
		}

		res, err := endpoint.HandleRequest(ctx, decoded)
		if err != nil {
			log.Errorf("http.MakeHandler: %s", err)
			renderer.RenderError(w, r, err)
			return
		}

		encoded, err := endpoint.EncodeResponse(ctx, res)
		if err != nil {
			log.Errorf("http.MakeHandler: %s", err)
			renderer.RenderError(w, r, err)
			return
		}

		if err = renderer.Render(w, r, encoded); err != nil {
			log.Errorf("http.MakeHandler: %s", err)
			renderer.RenderError(w, r, err)
		}
	}
}
