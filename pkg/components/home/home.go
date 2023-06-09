package home

import "github.com/itsnoproblem/prmry/pkg/components"

type HomeView struct {
	Providers    []string
	ProvidersMap map[string]string
	components.BaseComponent
}

type LoginView struct {
	Providers    []string
	ProvidersMap map[string]string
	components.BaseComponent
}
