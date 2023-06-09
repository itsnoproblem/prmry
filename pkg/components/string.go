package components

import (
	"strings"
)

func TrimWordsToMaxCharacters(maxChars int, text string) string {
	var (
		trimmed = ""
		strlen  = 0
	)

	fields := strings.Fields(text)
	for i, f := range fields {
		strlen += len(f) + 1
		if strlen > maxChars {
			trimmed = strings.TrimSuffix(trimmed, " ") + "..."
			return trimmed
		}

		if i < len(fields) {
			trimmed += f + " "
		}
	}

	return trimmed
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
