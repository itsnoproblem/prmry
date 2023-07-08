package redirect

import "github.com/itsnoproblem/prmry/internal/components"

type View struct {
	Status   int
	Location string
	components.BaseComponent
}
