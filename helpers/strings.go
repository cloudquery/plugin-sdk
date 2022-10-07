package helpers

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cast"
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

// NormalizeNewlines normalizes \r\n (windows) and \r (mac)
// into \n (unix)
func NormalizeNewlines(d []byte) []byte {
	// replace CR LF \r\n (windows) with LF \n (unix)
	d = bytes.ReplaceAll(d, []byte{13, 10}, []byte{10})
	// replace CF \r (mac) with LF \n (unix)
	d = bytes.ReplaceAll(d, []byte{13}, []byte{10})
	return d
}

func ToStringSliceE(i interface{}) ([]string, error) {
	switch v := i.(type) {
	case *[]string:
		return cast.ToStringSliceE(*v)
	default:
		return cast.ToStringSliceE(i)
	}
}
