package flow

import (
	"fmt"
	"strings"
)

const (
	ConditionTypeEquals     ConditionType = "equals"
	ConditionTypeContains   ConditionType = "contains"
	ConditionTypeStartsWith ConditionType = "starts with"
	ConditionTypeEndsWith   ConditionType = "ends with"

	FieldTypePrompt FieldType = "prompt"
)

type ConditionType string
type FieldType string

func (t ConditionType) String() string {
	return string(t)
}

func (t FieldType) String() string {
	return string(t)
}

type Flow struct {
	ID         string
	Name       string
	Conditions []Condition
	RequireAll bool
	Response   string
}

func (f Flow) Respond(prompt string) (string, error) {
	for _, cond := range f.Conditions {
		matches, err := cond.Matches(prompt)
		if err != nil {
			return "", fmt.Errorf("Flow.Respond: %s", err)
		}

		if matches {
			return f.Response, nil
		}
	}

	return "", nil
}

type Rule struct {
	Conditions []Condition
	RequireAll bool
	Response   string
}

func (r Rule) Matches(prompt string) (bool, error) {
	var isMatch bool
	var err error

	for _, cond := range r.Conditions {
		isMatch, err = cond.Matches(prompt)
		if err != nil {
			return false, fmt.Errorf("Rule.Matches: %s", err)
		}

		if (!r.RequireAll && isMatch) || (r.RequireAll && !isMatch) {
			return isMatch, nil
		}
	}

	return isMatch, nil
}

type Condition struct {
	Type  ConditionType
	Field FieldType
	Value string
}

func NewPromptCondition(t ConditionType, value string) Condition {
	return Condition{
		Type:  t,
		Field: FieldTypePrompt,
		Value: value,
	}
}

func (c Condition) Matches(field string) (bool, error) {
	fld := strings.ToLower(field)
	val := strings.ToLower(c.Value)

	switch c.Type {
	case ConditionTypeContains:
		return strings.Contains(fld, val), nil
	case ConditionTypeEquals:
		return fld == val, nil
	case ConditionTypeStartsWith:
		return strings.HasPrefix(fld, val), nil
	case ConditionTypeEndsWith:
		return strings.HasSuffix(fld, val), nil
	}

	return false, fmt.Errorf("Condition.Matches: unknown condition type: %s", c.Type)
}

func SupportedConditions() []string {
	return []string{
		ConditionTypeContains.String(),
		ConditionTypeEquals.String(),
		ConditionTypeStartsWith.String(),
		ConditionTypeEndsWith.String(),
	}
}

func SupportedFields() []string {
	return []string{
		FieldTypePrompt.String(),
	}
}
