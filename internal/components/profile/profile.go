package profile

import "github.com/itsnoproblem/prmry/internal/components"

type APIKeyView struct {
	Name      string
	Key       string
	CreatedAt string
}

type ProfileView struct {
	Provider  string
	Email     string
	Name      string
	AvatarURL string
	APIKeys   []APIKeyView
	components.BaseComponent
}
