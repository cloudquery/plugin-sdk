package remoteoauth

import (
	"golang.org/x/oauth2"
)

// NewTokenSource creates a new token source.
// Deprecated: Use oauth2.StaticTokenSource directly instead.
func NewTokenSource(opts ...TokenSourceOption) (oauth2.TokenSource, error) {
	t := &tokenSource{}
	for _, opt := range opts {
		opt(t)
	}
	return oauth2.StaticTokenSource(&t.currentToken), nil
}

type tokenSource struct {
	currentToken oauth2.Token
}
