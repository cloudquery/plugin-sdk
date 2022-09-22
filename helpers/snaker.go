package helpers

import (
	"strings"
	"unicode"
)

func AddInitialism(i string) {
	commonInitialisms[i] = struct{}{}
}

// ToSnake converts a given string to snake case
func ToSnake(s string) string {
	var result string
	var words []string
	var lastPos int
	rs := []rune(s)

	for i := 0; i < len(rs); i++ {
		if i > 0 && unicode.IsUpper(rs[i]) {
			if initialism := startsWithInitialism(s[lastPos:]); initialism != "" {
				words = append(words, initialism)

				i = lastPos + len(initialism)
				lastPos = i
				continue
			}

			words = append(words, s[lastPos:i])
			lastPos = i
		}
	}

	// append the last word
	if s[lastPos:] != "" {
		//handle plurals of initialisms like CDNs, ARNs, IDs
		if w := s[lastPos:]; w == "s" {
			words[len(words)-1] = words[len(words)-1] + w
		} else {
			words = append(words, s[lastPos:])
		}
	}

	for k, word := range words {
		if k > 0 {
			result += "_"
		}

		result += strings.ToLower(word)
	}

	return result
}

// ToCamel returns a string converted from snake case to uppercase
func ToCamel(s string) string {
	var result string

	words := strings.Split(s, "_")
	for i, word := range words {
		if exception := snakeToCamelExceptions[word]; len(exception) > 0 {
			result += exception
			continue
		}

		if i > 0 {
			upper := strings.ToUpper(word)
			if _, ok := commonInitialisms[upper]; len(s) > i-1 && ok {
				result += upper
				continue
			}
		}

		if (i > 0) && len(word) > 0 {
			w := []rune(word)
			w[0] = unicode.ToUpper(w[0])
			result += string(w)
		} else {
			result += word
		}
	}

	return result
}

// startsWithInitialism returns the initialism if the given string begins with it
func startsWithInitialism(s string) string {
	var initialism string
	// the longest initialism is 5 char, the shortest 2
	for i := 1; i <= len(s) && i <= 5; i++ {
		if _, ok := commonInitialisms[s[:i]]; len(s) > i-1 && ok {
			initialism = s[:i]
		}
	}
	return initialism
}

// commonInitialisms taken from https://github.com/golang/lint/blob/master/lint.go
var commonInitialisms = map[string]struct{}{
	"ACL":   {},
	"API":   {},
	"ASCII": {},
	"CPU":   {},
	"CSS":   {},
	"DNS":   {},
	"EOF":   {},
	"GUID":  {},
	"HTML":  {},
	"HTTP":  {},
	"HTTPS": {},
	"ID":    {},
	"IP":    {},
	"JSON":  {},
	"LHS":   {},
	"QPS":   {},
	"RAM":   {},
	"RHS":   {},
	"RPC":   {},
	"SLA":   {},
	"SMTP":  {},
	"SQL":   {},
	"SSH":   {},
	"TCP":   {},
	"TLS":   {},
	"TTL":   {},
	"UDP":   {},
	"UI":    {},
	"UID":   {},
	"UUID":  {},
	"URI":   {},
	"URL":   {},
	"UTF8":  {},
	"VM":    {},
	"XML":   {},
	"XMPP":  {},
	"XSRF":  {},
	"XSS":   {},
}

// add exceptions here for things that are not automatically convertable
var snakeToCamelExceptions = map[string]string{
	"oauth": "OAuth",
}
