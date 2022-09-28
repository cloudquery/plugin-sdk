package caser

import (
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Caser struct {
	initialisms            map[string]bool
	camelToSnakeExceptions map[string]string
	snakeToCamelException  map[string]string
}

type Option func(*Caser)

// WithCustomInitialims allows to specify custom initialisms for caser.
func WithCustomInitialims(fields map[string]bool) Option {
	return func(c *Caser) {
		for k, v := range fields {
			c.initialisms[k] = v
		}
	}
}

// WithCustomExceptions allows to specify custom exceptions for caser.
// The parameter is a map of camel:snake values like map[string]string{"oauth":"OAuth"}
func WithCustomExceptions(fields map[string]string) Option {
	return func(c *Caser) {
		for k, v := range fields {
			c.camelToSnakeExceptions[v] = k
			c.snakeToCamelException[k] = v
		}
	}
}

func reverseMap(m map[string]string) map[string]string {
	n := make(map[string]string, len(m))
	for k, v := range m {
		n[v] = k
	}
	return n
}

// New creates a new instance of caser
func New(opts ...Option) *Caser {
	c := &Caser{
		initialisms:            make(map[string]bool),
		camelToSnakeExceptions: make(map[string]string),
		snakeToCamelException:  make(map[string]string),
	}
	for k, v := range commonInitialisms {
		c.initialisms[k] = v
	}
	for k, v := range commonExceptions {
		c.snakeToCamelException[k] = v
		c.camelToSnakeExceptions[v] = k
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// ToSnake converts a given string to snake case
func (c *Caser) ToSnake(s string) string {
	if s == "" {
		return s
	}
	var result string
	var words []string
	var lastPos int
	rs := []rune(s)

	for i := 0; i < len(rs); i++ {
		if i > 0 && unicode.IsUpper(rs[i]) {
			if initialism := c.startsWithInitialism(s[lastPos:]); initialism != "" {
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

		if exception, ok := c.camelToSnakeExceptions[word]; ok {
			result += exception
			continue
		}

		result += strings.ToLower(word)
	}

	return result
}

// ToPascal returns a string converted from snake case to pascal case
func (c *Caser) ToPascal(s string) string {
	if s == "" {
		return s
	}
	result := c.ToCamel(s)
	csr := cases.Title(language.Und, cases.NoLower)
	return csr.String(result)
}

// ToCamel returns a string converted from snake case to camel case
func (c *Caser) ToCamel(s string) string {
	if s == "" {
		return s
	}
	var result string

	words := strings.Split(s, "_")
	for i, word := range words {
		if exception, ok := c.snakeToCamelException[word]; ok {
			result += exception
			continue
		}

		if i > 0 {
			upper := strings.ToUpper(word)
			if len(s) > i-1 && c.initialisms[upper] {
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
func (c *Caser) startsWithInitialism(s string) string {
	var initialism string
	// the longest initialism is 5 char, the shortest 2
	// we choose the longest match
	for i := 1; i <= len(s) && i <= 5; i++ {
		if len(s) > i-1 && c.initialisms[s[:i]] && len(s[:i]) > len(initialism) {
			initialism = s[:i]
		}
	}
	return initialism
}
