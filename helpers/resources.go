package helpers

import (
	"fmt"
	"strings"
)

func FormatSlice(a []string) string {
	q := make([]string, len(a))
	for i, s := range a {
		q[i] = fmt.Sprintf("%q", s)
	}
	return fmt.Sprintf("[\n\t%s\n]", strings.Join(q, ",\n\t"))
}
