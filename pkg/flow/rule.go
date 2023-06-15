package flow

import (
	"fmt"
	"strings"
)

const (
	ConditionTypeEquals      ConditionType = "equals"
	ConditionTypeNotEquals   ConditionType = "does not equal"
	ConditionTypeContains    ConditionType = "contains"
	ConditionTypeNotContains ConditionType = "does not contain"
	ConditionTypeStartsWith  ConditionType = "starts with"
	ConditionTypeEndsWith    ConditionType = "ends with"

	FieldSourceInput FieldSource = "input message"
	FieldSourceFlow  FieldSource = "interaction result from another flow"
)

type Rule struct {
	Field     FieldSource
	Condition ConditionType
	Value     string
}

func (c Rule) Matches(field string) (bool, error) {
	fld := strings.ToLower(field)
	val := strings.ToLower(c.Value)

	switch c.Condition {
	case ConditionTypeContains:
		return strings.Contains(fld, val), nil
	case ConditionTypeEquals:
		return fld == val, nil
	case ConditionTypeStartsWith:
		return strings.HasPrefix(fld, val), nil
	case ConditionTypeEndsWith:
		return strings.HasSuffix(fld, val), nil
	}

	return false, fmt.Errorf("Rule.Matches: unknown condition type: %s", c.Condition)
}

type ConditionType string

func (t ConditionType) String() string {
	return string(t)
}

type FieldSource string

func (s FieldSource) String() string {
	return string(s)
}

func SupportedConditions() []string {
	return []string{
		ConditionTypeContains.String(),
		ConditionTypeNotContains.String(),
		ConditionTypeEquals.String(),
		ConditionTypeNotEquals.String(),
		ConditionTypeStartsWith.String(),
		ConditionTypeEndsWith.String(),
	}
}

func SupportedFields() []string {
	return []string{
		FieldSourceInput.String(),
		FieldSourceFlow.String(),
	}
}
