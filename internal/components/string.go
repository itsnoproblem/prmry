package components

import (
	"sort"
)

type SortedMap map[string]string

func (s SortedMap) Keys() []string {
	keys := make([]string, 0)
	for k, _ := range s {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func SelectedIfTrue(cond bool) string {
	if cond {
		return "selected"
	}
	return ""
}

func TrueFalse(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
