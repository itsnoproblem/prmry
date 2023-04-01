package htmx

import "net/http"

const (
	HeaderHXRequest     = "HX-Request"
	HeaderHXRedirect    = "HX-Redirect"
	HeaderHXTarget      = "HX-Target"
	HeaderHXTriggerName = "HX-Trigger-Name"
	HeaderHXCurrentURL  = "HX-Current-URL"
)

// IsHXRequest returns true if the HX-Request header is "true"
func IsHXRequest(r *http.Request) bool {
	return r.Header.Get(HeaderHXRequest) == "true"
}

// Redirect sets the HX-Redirect header to the value of url
func Redirect(w http.ResponseWriter, url string) {
	w.Header().Set(HeaderHXRedirect, url)
}
