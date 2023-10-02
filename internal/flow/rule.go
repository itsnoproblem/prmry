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
)

type ConditionType string

func (t ConditionType) String() string {
	return string(t)
}

type Rule struct {
	Field     Field
	Condition ConditionType
	Value     string
}

func (c Rule) Matches(field string) (bool, error) {
	matches := false
	fld := strings.ToLower(field)
	val := strings.ToLower(c.Value)

	for strings.HasPrefix(val, "\n") {
		val = strings.TrimPrefix(val, "\n")
	}
	for strings.HasPrefix(fld, "\n") {
		fld = strings.TrimPrefix(fld, "\n")
	}

	switch c.Condition {
	case ConditionTypeContains:
		matches = strings.Contains(fld, val)
		break
	case ConditionTypeNotContains:
		matches = !strings.Contains(fld, val)
		break
	case ConditionTypeEquals:
		matches = fld == val
		break
	case ConditionTypeNotEquals:
		matches = fld != val
		break
	case ConditionTypeStartsWith:
		matches = strings.HasPrefix(fld, val)
		break
	case ConditionTypeEndsWith:
		matches = strings.HasSuffix(fld, val)
		break
	default:
		return false, fmt.Errorf("Rule.Matches: unknown condition type: %s", c.Condition)
	}

	return matches, nil
}

func SupportedConditions() map[string]string {
	return map[string]string{
		ConditionTypeContains.String():    ConditionTypeContains.String(),
		ConditionTypeNotContains.String(): ConditionTypeNotContains.String(),
		ConditionTypeEquals.String():      ConditionTypeEquals.String(),
		ConditionTypeNotEquals.String():   ConditionTypeNotEquals.String(),
		ConditionTypeStartsWith.String():  ConditionTypeStartsWith.String(),
		ConditionTypeEndsWith.String():    ConditionTypeEndsWith.String(),
	}
}
