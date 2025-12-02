package remoteoauth

import (
	"context"
	"time"

	"golang.org/x/oauth2"
)

type TokenSourceOption func(source *tokenSource)

// WithAccessToken sets the access token, token type and expiry time for the token source.
// Deprecated: Use WithToken instead.
func WithAccessToken(token, tokenType string, expiry time.Time) TokenSourceOption {
	return func(t *tokenSource) {
		t.currentToken = oauth2.Token{
			AccessToken: token,
			TokenType:   tokenType,
			Expiry:      expiry,
		}
	}
}

// WithToken sets the default token for the token source.
// Deprecated: Use oauth2.StaticTokenSource directly instead.
func WithToken(token oauth2.Token) TokenSourceOption {
	return func(t *tokenSource) {
		t.currentToken = token
	}
}

// WithDefaultContext sets the default context for the token source, used when creating a new token request.
// Deprecated: not used in the current implementation.
func WithDefaultContext(_ context.Context) TokenSourceOption {
	return func(t *tokenSource) {}
}
