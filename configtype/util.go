package configtype

import "strings"

func patternCases(cases ...string) string {
	return "(" + strings.Join(cases, "|") + ")"
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func plural(n int) string {
	if n != 1 {
		return "s"
	}
	return ""
}
