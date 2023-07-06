package api

import (
	"io"
	"net/http"

	"github.com/google/martian/log"
)

func MakeHandler(e Endpoint) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoded, err := e.DecodeRequest(r.Context(), r)
		if err != nil {
			handleError(err, w, r)
			return
		}

		res, err := e.HandleRequest(r.Context(), decoded)
		if err != nil {
			handleError(err, w, r)
			return
		}

		formatted, err := e.FormatResponse(res)
		if err != nil {
			handleError(err, w, r)
			return
		}

		encoded := SuccessResponse{
			Data: formatted,
		}

		if err = encoded.Render(w, r); err != nil {
			log.Errorf("api-server.MakeHandler: %s", err)
		}
	}
}

func handleError(err error, w http.ResponseWriter, r *http.Request) {
	res := ErrorInternal(err)
	if renderErr := res.Render(w, r); renderErr != nil {
		log.Errorf("api-server.handleError: %s", err)
	}
}

type HTMXRenderer interface {
	Render(w io.Writer, data interface{}) error
	RenderError(w io.Writer, err error) error
}
