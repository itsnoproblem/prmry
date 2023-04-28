package auth

import (
	"context"
	"encoding/gob"
	"errors"
	"log"
	"net/http"
	"strings"
)

func Middleware(secret Byte32) func(http.Handler) http.Handler {
	f := func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var usr User

			if !strings.HasPrefix(r.URL.Path, "/auth") {
				gobEncodedValue, err := ReadEncrypted(r, CookieName, secret)
				if err != nil {
					switch {
					case errors.Is(err, http.ErrNoCookie):
						cookie, err := NewCookie(CookieName, User{})
						if err != nil {
							w.Write([]byte(err.Error()))
							http.Error(w, err.Error(), http.StatusInternalServerError)
						}
						gobEncodedValue = cookie.Value

					case errors.Is(err, ErrInvalidValue):
						w.Write([]byte(err.Error()))
						http.Error(w, "invalid cookie", http.StatusBadRequest)
						return
					default:
						w.Write([]byte(err.Error()))
						http.Error(w, "server error: "+err.Error(), http.StatusInternalServerError)
						return
					}
				}

				// Create an strings.Reader containing the gob-encoded value.
				reader := strings.NewReader(gobEncodedValue)

				// Decode it into the User type. Notice that we need to pass a *pointer* to
				// the Decode() target here?
				if err := gob.NewDecoder(reader).Decode(&usr); err != nil {
					http.Error(w, "server error: "+err.Error(), http.StatusInternalServerError)
					return
				}
			}

			ctx := context.WithValue(r.Context(), ContextKey, usr)
			reqWithUser := r.WithContext(ctx)

			log.Printf("User: %v\n", usr)
			h.ServeHTTP(w, reqWithUser)
		}
		return http.HandlerFunc(fn)
	}

	return f
}
