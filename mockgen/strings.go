package mockgen

import (
	"reflect"
	"strings"
)

// MethodHasAnyPrefix returns true if the method name has any of the given prefixes.
func MethodHasAnyPrefix(m reflect.Method, prefixes []string) bool {
	for i := range prefixes {
		if strings.HasPrefix(m.Name, prefixes[i]) {
			return true
		}
	}
	return false
}

// MethodHasAnySuffix returns true if the method name has any of the given suffixes.
func MethodHasAnySuffix(m reflect.Method, suffixes []string) bool {
	for i := range suffixes {
		if strings.HasSuffix(m.Name, suffixes[i]) {
			return true
		}
	}
	return false
}
