package flow

type BySpecificity []Flow

func (a BySpecificity) Len() int      { return len(a) }
func (a BySpecificity) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a BySpecificity) Less(i, j int) bool {
	return rulesetSpecificity(a[i]) > rulesetSpecificity(a[j]) // Use > for descending order
}

func comparatorSpecificity(comp ConditionType) int {
	switch comp {
	case ConditionTypeEquals, ConditionTypeNotEquals:
		return 3 // most specific
	case ConditionTypeStartsWith, ConditionTypeEndsWith:
		return 2
	case ConditionTypeContains, ConditionTypeNotContains:
		return 1 // least specific
	}
	return 0
}

func rulesetSpecificity(flw Flow) int {
	specificity := 0
	for _, rule := range flw.Triggers {
		specificity += comparatorSpecificity(rule.Condition)
	}
	if flw.RequireAll {
		specificity += len(flw.Triggers) // Inclusive rulesets are more specific
	}
	return specificity
}
