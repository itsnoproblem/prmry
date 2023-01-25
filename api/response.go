package api

import (
	"github.com/go-chi/render"
	"net/http"
)

type SuccessResponse struct {
	Data interface{} `json:"data"`
}

func (res SuccessResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)

	return nil
}

type errorResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *errorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	render.JSON(w, r, map[string]interface{}{
		"error": e,
	})

	return nil
}

func ErrorInternal(err error) render.Renderer {
	return &errorResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		AppCode:        http.StatusInternalServerError,
		StatusText:     "Internal Server Error",
		ErrorText:      err.Error(),
	}
}
