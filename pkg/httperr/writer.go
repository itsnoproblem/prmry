package httperr

import "net/http"

// Unauthorized writes a 401 unauthorized response to w
func Unauthorized(msg string, err error, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`<div class="error">` + msg + `</div>`))
}

// BadRequest writes a 400 Bad Request response to w
func BadRequest(msg string, err error, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(`<div class="error">` + msg + `</div>`))
}

// Internal writes a 500 internal error response to w
func Internal(msg string, err error, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`<div class="error">` + msg + `</div>`))
}
