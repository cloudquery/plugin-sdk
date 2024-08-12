package remoteoauth

import (
	"context"
	"time"

	"golang.org/x/oauth2"
)

type TokenSourceOption func(source *tokenSource)

func WithAccessToken(token, tokenType string, expiry time.Time) TokenSourceOption {
	return func(t *tokenSource) {
		t.currentToken = oauth2.Token{
			AccessToken: token,
			TokenType:   tokenType,
			Expiry:      expiry,
		}
	}
}

func WithDefaultContext(ctx context.Context) TokenSourceOption {
	return func(t *tokenSource) {
		t.defaultContext = ctx
	}
}

func withNoWrap() TokenSourceOption {
	return func(t *tokenSource) {
		t.noWrap = true
	}
}
