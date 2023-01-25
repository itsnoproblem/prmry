package api

import (
	"github.com/google/martian/log"
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

func handleError(err error, w http.ResponseWriter, r *http.Request) {
	res := ErrorInternal(err)
	if renderErr := res.Render(w, r); renderErr != nil {
		log.Errorf("api.handleError: %s", err)
	}
}
