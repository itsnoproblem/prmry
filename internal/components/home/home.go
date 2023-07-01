package home

import "github.com/itsnoproblem/prmry/internal/components"

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
