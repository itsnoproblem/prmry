package authenticating

import (
	"context"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth/gothic"

	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/itsnoproblem/prmry/internal/components"
	"github.com/itsnoproblem/prmry/internal/components/home"
	"github.com/itsnoproblem/prmry/internal/htmx"
)

type AuthService interface {
	CreateUser(ctx context.Context, usr auth.User) (id string, err error)
	DeleteUser(ctx context.Context, id string) error
	SaveUserWithOAuthConnection(ctx context.Context, usr auth.User, provider, providerUserID string) error
	GetUserByProvider(ctx context.Context, provider, providerUserID string) (usr auth.User, exists bool, err error)
	GetUserByEmail(ctx context.Context, email string) (usr auth.User, exists bool, err error)
}

type Renderer interface {
	Render(w http.ResponseWriter, r *http.Request, cmp components.Component) error
	RenderTemplComponent(w http.ResponseWriter, r *http.Request, fullPage, fragment templ.Component) error
	RenderError(w http.ResponseWriter, r *http.Request, err error)
}

type Resource struct {
	authService AuthService
	renderer    Renderer
	secret      auth.Byte32
}

func NewResource(renderer Renderer, authSecret auth.Byte32, authService AuthService) (Resource, error) {
	gothic.GetProviderName = func(req *http.Request) (string, error) {
		return chi.URLParam(req, "provider"), nil
	}

	return Resource{
		authService: authService,
		renderer:    renderer,
		secret:      authSecret,
	}, nil
}

func (rs Resource) RouteHandler() func(chi.Router) {
	return func(r chi.Router) {
		r.Get("/auth/{provider}", rs.AuthHandler)
		r.Get("/auth/{provider}/callback", rs.AuthSuccessHandler)
		r.Get("/auth/logout", rs.LogoutHandler)
		r.Get("/auth/logout/{provider}", rs.LogoutHandler)
	}
}

func (rs Resource) AuthHandler(w http.ResponseWriter, r *http.Request) {
	// try to get the user without re-authenticating
	if usr, err := gothic.CompleteUserAuth(w, r); err == nil {
		authUser, exists, err := rs.authService.GetUserByProvider(r.Context(), usr.Provider, usr.UserID)
		if err != nil {
			rs.renderer.RenderError(w, r, err)
			return
		}

		if !exists {
			authUser = auth.User{
				Name:  usr.FirstName + " " + usr.LastName,
				Email: usr.Email,
			}
			if err = rs.authService.SaveUserWithOAuthConnection(r.Context(), authUser, usr.Provider, usr.UserID); err != nil {
				rs.renderer.RenderError(w, r, err)
				return
			}
		}

		if err := rs.saveUserCookie(w, authUser); err != nil {
			rs.renderer.RenderError(w, r, err)
			return
		}

		cmp := home.HomeView{}
		if err = rs.renderer.RenderTemplComponent(w, r, home.HomePage(cmp), home.HomeFragment(cmp)); err != nil {
			rs.renderer.RenderError(w, r, err)
			return
		}
	} else {
		url, err := gothic.GetAuthURL(w, r)
		if err != nil {
			rs.renderer.RenderError(w, r, err)
			return
		}

		if htmx.IsHXRequest(r.Context()) {
			w.Header().Set("HX-Redirect", url)
		} else {
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
			return
		}
	}
}

func (rs Resource) AuthSuccessHandler(w http.ResponseWriter, r *http.Request) {
	finishWithRedirect := func() {
		if htmx.IsHXRequest(r.Context()) {
			htmx.Redirect(w, "/")
			return
		}

		if r != nil && r.URL.Path != "/" {
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}

	userFromProvider, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}

	if userFromProvider.Email == "" {
		rs.renderer.RenderError(w, r, fmt.Errorf("user from provider missing email"))
		return
	}

	userByEmail, userByEmailExists, err := rs.authService.GetUserByEmail(r.Context(), userFromProvider.Email)
	if err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}

	if !userByEmailExists {
		name := userFromProvider.NickName
		if userFromProvider.FirstName != "" {
			name = userFromProvider.FirstName + " " + userFromProvider.LastName
		}
		userByEmail = auth.User{
			Name:      name,
			Email:     userFromProvider.Email,
			Nickname:  userFromProvider.NickName,
			AvatarURL: userFromProvider.AvatarURL,
		}

		userByEmail.ID, err = rs.authService.CreateUser(r.Context(), userByEmail)
		if err != nil {
			rs.renderer.RenderError(w, r, err)
			return
		}

		if err = rs.authService.SaveUserWithOAuthConnection(r.Context(), userByEmail, userFromProvider.Provider, userFromProvider.UserID); err != nil {
			rs.renderer.RenderError(w, r, err)
			return
		}
	}

	if err := rs.saveUserCookie(w, userByEmail); err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}

	finishWithRedirect()
}

func (rs Resource) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	err := gothic.Logout(w, r)
	if err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}

	if err = rs.deleteUserCookie(w); err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}

	if htmx.IsHXRequest(r.Context()) {
		w.Header().Add("HX-Redirect", "/")
	}

	if r != nil && r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (rs Resource) saveUserCookie(w http.ResponseWriter, user auth.User) error {
	cookie, err := auth.NewCookie(auth.CookieName, user)
	if err != nil {
		return fmt.Errorf("saveUserCookie: %s", err)
	}

	if err = auth.WriteEncrypted(w, cookie, rs.secret); err != nil {
		return fmt.Errorf("saveUserCookie: %s", err)
	}

	return nil
}

func (rs Resource) deleteUserCookie(w http.ResponseWriter) error {
	if err := rs.saveUserCookie(w, auth.User{}); err != nil {
		return fmt.Errorf("deleteUserCookie: %s", err)
	}

	return nil
}
