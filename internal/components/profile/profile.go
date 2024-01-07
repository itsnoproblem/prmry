package profile

import "github.com/itsnoproblem/prmry/internal/components"

type APIKeyView struct {
	Name      string
	Key       string
	CreatedAt string
	components.BaseComponent
}

type ProfileView struct {
	Provider  string
	Email     string
	Name      string
	AvatarURL string
	APIKeys   []APIKeyView
	components.BaseComponent
}

type APIKeySuccessView struct {
	Name string
	components.BaseComponent
}
