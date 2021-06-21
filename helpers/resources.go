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
