package api

import (
	"github.com/google/martian/log"
	"io"
	"net/http"
)

func MakeHandler(e Endpoint) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoded, err := e.Decode(r.Context(), r)
		if err != nil {
			handleError(err, w, r)
			return
		}

		res, err := e.Handle(r.Context(), decoded)
		if err != nil {
			handleError(err, w, r)
			return
		}

		formatted, err := e.Format(res)
		if err != nil {
			handleError(err, w, r)
			return
		}

		encoded := SuccessResponse{
			Data: formatted,
		}

		if err = encoded.Render(w, r); err != nil {
			log.Errorf("api.MakeHandler: %s", err)
		}
	}
}

func MakeHTMXHandler(e Endpoint, renderer HTMXRenderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoded, err := e.Decode(r.Context(), r)
		if err != nil {
			handleError(err, w, r)
			return
		}

		res, err := e.Handle(r.Context(), decoded)
		if err != nil {
			handleErrorHTMX(err, w, renderer)
			return
		}

		//formatted, err := e.Format(res)
		//if err != nil {
		//	handleErrorHTMX(err, w, renderer)
		//	return
		//}

		if err = renderer.Render(w, res); err != nil {
			log.Errorf("api.MakeHandler: %s", err)
		}
	}
}

func handleError(err error, w http.ResponseWriter, r *http.Request) {
	res := ErrorInternal(err)
	if renderErr := res.Render(w, r); renderErr != nil {
		log.Errorf("api.handleError: %s", err)
	}
}

type HTMXRenderer interface {
	Render(w io.Writer, data interface{}) error
	RenderError(w io.Writer, err error) error
}

func handleErrorHTMX(err error, w http.ResponseWriter, r HTMXRenderer) {
	r.RenderError(w, err)
}
