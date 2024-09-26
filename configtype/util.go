package configtype

import "strings"

func patternCases(cases ...string) string {
	return "(" + strings.Join(cases, "|") + ")"
}
