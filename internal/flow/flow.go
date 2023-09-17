package flow

import (
	"time"
)

const (
	ParamTypeString ParamType = "string"
	ParamTypeNumber ParamType = "number"
)

type Flow struct {
	ID          string
	UserID      string
	Name        string
	Rules       []Rule
	RequireAll  bool
	Prompt      string
	PromptArgs  []Field
	CreatedAt   time.Time
	UpdatedAt   time.Time
	InputParams []InputParam
}

type ParamType string
type InputParam struct {
	Type       ParamType
	Key        string
	IsRequired bool
}
