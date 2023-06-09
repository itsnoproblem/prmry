package profile

import "github.com/itsnoproblem/prmry/pkg/components"

type ProfileView struct {
	Provider  string
	Email     string
	Name      string
	AvatarURL string
	components.BaseComponent
}
