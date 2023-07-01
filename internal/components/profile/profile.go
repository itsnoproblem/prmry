package profile

import "github.com/itsnoproblem/prmry/internal/components"

type ProfileView struct {
	Provider  string
	Email     string
	Name      string
	AvatarURL string
	components.BaseComponent
}
