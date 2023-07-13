package htmx

import (
	"context"
	"net/http"
)

const (
	ContextKeyHxRequest     = "HXRequest"
	ContextKeyHXTarget      = "HXTarget"
	ContextKeyHXTriggerName = "HXTriggerName"
	ContextKeyHXCurrentURL  = "HXCurrentURL"
)

// IsHXRequest returns true if the HX-Request header is "true"
func IsHXRequest(ctx context.Context) bool {
	return ctx.Value(ContextKeyHxRequest) == "true"
}

func GetTarget(ctx context.Context) string {
	return ctx.Value(ContextKeyHXTarget).(string)
}

func GetTriggerName(ctx context.Context) string {
	return ctx.Value(ContextKeyHXTriggerName).(string)
}

func CurrentURL(ctx context.Context) string {
	return ctx.Value(ContextKeyHXCurrentURL).(string)
}

func Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextKeyHxRequest, r.Header.Get(HeaderHXRequest))
		ctx = context.WithValue(ctx, ContextKeyHXTarget, r.Header.Get(HeaderHXTarget))
		ctx = context.WithValue(ctx, ContextKeyHXTriggerName, r.Header.Get(HeaderHXTriggerName))
		ctx = context.WithValue(ctx, ContextKeyHXCurrentURL, r.Header.Get(HeaderHXCurrentURL))

		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
