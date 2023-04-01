package authorizing

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/auth"
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/cookies"
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/htmx"
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/httperr"
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/user"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"log"
	"net/http"
)

const (
	paramNameProvider = "provider"
)

type AuthService interface {
	CreateUser(ctx context.Context, usr user.User) (id string, err error)
	DeleteUser(ctx context.Context, id string) error
	SaveUserWithOAuthConnection(ctx context.Context, usr user.User, provider, providerUserID string) error
	GetUserByProvider(ctx context.Context, provider, providerUserID string) (usr user.User, exists bool, err error)
	GetUserByEmail(ctx context.Context, email string) (usr user.User, exists bool, err error)
}

type Resource struct {
	authService AuthService
	renderer    Renderer
	secret      auth.Byte32
}

func NewResource(authSecret auth.Byte32, authService AuthService) (Resource, error) {
	renderer, err := NewRenderer()
	if err != nil {
		return Resource{}, fmt.Errorf("authorizing.NewResource: %s", err)
	}

	gothic.GetProviderName = func(req *http.Request) (string, error) {
		return chi.URLParam(req, paramNameProvider), nil
	}

	return Resource{
		authService: authService,
		renderer:    *renderer,
		secret:      authSecret,
	}, nil
}

func (rs Resource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get(fmt.Sprintf("/{%s}", paramNameProvider), rs.AuthHandler)
	r.Get(fmt.Sprintf("/{%s}/callback", paramNameProvider), rs.AuthSuccessHandler)
	r.Get(fmt.Sprintf("/logout/{%s}", paramNameProvider), rs.LogoutHandler)
	return r
}

func (rs Resource) AuthHandler(w http.ResponseWriter, r *http.Request) {
	// try to get the user without re-authenticating
	if user, err := gothic.CompleteUserAuth(w, r); err == nil {
		if err := rs.saveUserCookie(w, user); err != nil {
			httperr.Internal(err.Error(), err, w, r)
			return
		}

		if err = rs.renderer.RenderLoginSuccess(w, NewUserView(user)); err != nil {
			httperr.Internal("AuthHandler: "+err.Error(), err, w, r)
			return
		}
	} else {
		url, err := gothic.GetAuthURL(w, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			httperr.BadRequest("AuthHandler: "+err.Error(), err, w, r)
			return
		}

		if htmx.IsHXRequest(r) {
			w.Header().Set("HX-Redirect", url)
		} else {
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
			return
		}
	}
}

func (rs Resource) AuthSuccessHandler(w http.ResponseWriter, r *http.Request) {
	finishWithRedirect := func() {
		if htmx.IsHXRequest(r) {
			htmx.Redirect(w, "/")
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}

	userFromProvider, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		httperr.Internal("CompleteUserAuth: "+err.Error(), err, w, r)
	}

	if userFromProvider.Email != "" {
		userByEmail, userByEmailExists, err := rs.authService.GetUserByEmail(r.Context(), userFromProvider.Email)
		if err != nil {
			httperr.Internal("AuthSuccessHandler: "+err.Error(), err, w, r)
		}

		if userByEmailExists {
			if err = rs.authService.SaveUserWithOAuthConnection(r.Context(), userByEmail, userFromProvider.Provider, userFromProvider.UserID); err != nil {
				httperr.Internal("failed to save oauth connection", err, w, r)
				return
			}

			finishWithRedirect()
			return
		}
	}

	existingUser, exists, err := rs.authService.GetUserByProvider(r.Context(), userFromProvider.Provider, userFromProvider.UserID)
	if err != nil {
		httperr.Internal("AuthSuccessHandler: "+err.Error(), err, w, r)
		return
	}

	if exists {
		log.Printf("Wuthenticated User: %s", existingUser.ID)
		if err = rs.authService.SaveUserWithOAuthConnection(r.Context(), existingUser, userFromProvider.Provider, userFromProvider.UserID); err != nil {
			httperr.Internal("failed to save oauth connection", err, w, r)
			return
		}
	} else {
		name := userFromProvider.NickName
		if userFromProvider.FirstName != "" {
			name = userFromProvider.Name
		}
		usr := user.User{
			Name:  name,
			Email: userFromProvider.Email,
		}

		usr.ID, err = rs.authService.CreateUser(r.Context(), usr)
		if err != nil {
			httperr.Internal("Failed to create user", err, w, r)
		}

		log.Printf("Created User: %s", usr.ID)

		if err = rs.authService.SaveUserWithOAuthConnection(r.Context(), usr, userFromProvider.Provider, userFromProvider.UserID); err != nil {
			log.Printf("ERROR (rollback): Failed to save oauth connection: %s", err)
			if err := rs.authService.DeleteUser(r.Context(), usr.ID); err != nil {
				log.Printf("ERROR: Failed to delete user during rollback: %s", usr.ID)
			}

			httperr.Internal("failed to save oauth connection", err, w, r)
			return
		}
	}

	if err := rs.saveUserCookie(w, userFromProvider); err != nil {
		httperr.Internal("saveUserCookie: "+err.Error(), err, w, r)
		return
	}

	finishWithRedirect()

}

func (rs Resource) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	err := gothic.Logout(w, r)
	if err != nil {
		httperr.Internal("LogoutHandler: "+err.Error(), err, w, r)
	}

	if err = rs.deleteUserCookie(w); err != nil {
		httperr.Internal("LogoutHandler: "+err.Error(), err, w, r)
	}

	w.Header().Add("HX-Redirect", "/")
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (rs Resource) saveUserCookie(w http.ResponseWriter, user goth.User) error {
	cookie, err := cookies.New(auth.CookieName, user)
	if err != nil {
		return fmt.Errorf("saveUserCookie: %s", err)
	}

	if err = cookies.WriteEncrypted(w, cookie, rs.secret); err != nil {
		return fmt.Errorf("saveUserCookie: %s", err)
	}

	return nil
}

func (rs Resource) deleteUserCookie(w http.ResponseWriter) error {
	if err := rs.saveUserCookie(w, goth.User{}); err != nil {
		return fmt.Errorf("deleteUserCookie: %s", err)
	}

	return nil
}
