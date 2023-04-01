package profiling

import (
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/htmx"
)

type HomeView struct {
	Providers    []string
	ProvidersMap map[string]string
	htmx.BaseComponent
}

type LoginView struct {
	Providers    []string
	ProvidersMap map[string]string
	htmx.BaseComponent
}
