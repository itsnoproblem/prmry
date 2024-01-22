package funnel

import "github.com/itsnoproblem/prmry/internal/components"

type FunnelFormView struct {
	components.BaseComponent
	ID   string
	Name string
	Path string
}
