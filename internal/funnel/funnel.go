package funnel

import (
	"github.com/itsnoproblem/prmry/internal/flow"
	"time"
)

type Funnel struct {
	UserID    string
	ID        string
	Name      string
	Path      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type WithFlows struct {
	Funnel
	Flows []flow.Flow
}

type Summary struct {
	ID        string
	Name      string
	Path      string
	FlowCount int
	UpdatedAt time.Time
}
