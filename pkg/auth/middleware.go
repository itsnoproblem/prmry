package auth

import (
	"context"
	"encoding/gob"
	"errors"
	"net/http"
	"strings"
)

func Middleware(secret Byte32) func(http.Handler) http.Handler {
	f := func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var usr User

			if r.Header.Get("X-Forwarded-Proto") == "https" && r.Proto != "https" {
				http.Redirect(w, r, r.URL.String(), http.StatusFound)
			}

			if !strings.HasPrefix(r.URL.Path, "/auth") {
				gobEncodedValue, err := ReadEncrypted(r, CookieName, secret)
				if err != nil {
					switch {
					case errors.Is(err, http.ErrNoCookie):
						cookie, err := NewCookie(CookieName, User{})
						if err != nil {
							http.Error(w, "create cookie: "+err.Error(), http.StatusInternalServerError)
						}
						gobEncodedValue = cookie.Value

					case errors.Is(err, ErrInvalidValue):
						http.Error(w, "invalid cookie: "+err.Error(), http.StatusBadRequest)
						return
					default:
						http.Error(w, "server error: "+err.Error(), http.StatusInternalServerError)
						return
					}
				}

				reader := strings.NewReader(gobEncodedValue)
				if err := gob.NewDecoder(reader).Decode(&usr); err != nil {
					http.Error(w, "server error: "+err.Error(), http.StatusInternalServerError)
					return
				}
			}

			ctx := context.WithValue(r.Context(), ContextKey, usr)
			reqWithUser := r.WithContext(ctx)

			h.ServeHTTP(w, reqWithUser)
		}
		return http.HandlerFunc(fn)
	}

	return f
}
