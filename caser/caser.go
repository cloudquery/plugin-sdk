package caser

import (
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// ToSnake converts a given string to snake case
func ToSnake(s string) string {
	if s == "" {
		return s
	}
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

// ToPascal returns a string converted from snake case to pascal case
func ToPascal(s string) string {
	if s == "" {
		return s
	}
	result := ToCamel(s)
	c := cases.Title(language.Und, cases.NoLower)

	return c.String(result)
}

// ToCamel returns a string converted from snake case to camel case
func ToCamel(s string) string {
	if s == "" {
		return s
	}
	var result string

	words := strings.Split(s, "_")
	for i, word := range words {
		if exception := snakeToCamelExceptions[word]; len(exception) > 0 {
			result += exception
			continue
		}

		if i > 0 {
			upper := strings.ToUpper(word)
			if len(s) > i-1 && commonInitialisms[upper] {
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
