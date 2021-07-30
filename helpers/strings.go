package helpers

import (
	"fmt"
	"sort"
	"strings"
)

func FormatSlice(a []string) string {
	// sort slice for consistency
	sort.Strings(a)
	q := make([]string, len(a))
	for i, s := range a {
		q[i] = fmt.Sprintf("%q", s)
	}
	return fmt.Sprintf("[\n\t%s\n]", strings.Join(q, ",\n\t"))
}

func HasDuplicates(resources []string) bool {
	dups := make(map[string]bool, len(resources))
	for _, r := range resources {
		if _, ok := dups[r]; ok {
			return true
		}
		dups[r] = true
	}
	return false
}
