package auth

import (
	"context"
	"encoding/gob"
	"errors"
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/cookies"
	"github.com/markbates/goth"
	"net/http"
	"strings"
)

func Middleware(secret Byte32) func(http.Handler) http.Handler {
	f := func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var user goth.User

			if !strings.HasPrefix(r.URL.Path, "/auth") {
				gobEncodedValue, err := cookies.ReadEncrypted(r, CookieName, secret)
				if err != nil {
					switch {
					case errors.Is(err, http.ErrNoCookie):
						cookie, err := cookies.New(CookieName, goth.User{})
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
						}
						gobEncodedValue = cookie.Value

					case errors.Is(err, cookies.ErrInvalidValue):
						http.Error(w, "invalid cookie", http.StatusBadRequest)
						return
					default:
						http.Error(w, "server error: "+err.Error(), http.StatusInternalServerError)
						return
					}
				}

				// Create an strings.Reader containing the gob-encoded value.
				reader := strings.NewReader(gobEncodedValue)

				// Decode it into the User type. Notice that we need to pass a *pointer* to
				// the Decode() target here?
				if err := gob.NewDecoder(reader).Decode(&user); err != nil {
					http.Error(w, "server error", http.StatusInternalServerError)
					return
				}
			}

			ctx := context.WithValue(r.Context(), ContextKey, user)
			reqWithUser := r.WithContext(ctx)

			//log.Printf("User: %v\n", user)
			h.ServeHTTP(w, reqWithUser)
		}
		return http.HandlerFunc(fn)
	}

	return f
}
