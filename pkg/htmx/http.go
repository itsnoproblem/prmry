package htmx

import "net/http"

const (
	HeaderHXRequest     = "HX-Request"
	HeaderHXRedirect    = "HX-Redirect"
	HeaderHXTarget      = "HX-Target"
	HeaderHXTriggerName = "HX-Trigger-Name"
	HeaderHXCurrentURL  = "HX-Current-URL"
)

// Redirect sets the HX-Redirect header to the value of url
func Redirect(w http.ResponseWriter, url string) {
	w.Header().Set(HeaderHXRedirect, url)
}
