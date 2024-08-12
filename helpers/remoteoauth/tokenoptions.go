package remoteoauth

import (
	"time"

	"golang.org/x/oauth2"
)

type tokenSourceOption func(source *tokenSource)

func WithAccessToken(token, tokenType string, expiry time.Time) tokenSourceOption {
	return func(t *tokenSource) {
		t.currentToken = oauth2.Token{
			AccessToken: token,
			TokenType:   tokenType,
			Expiry:      expiry,
		}
	}
}

func withNoWrap() tokenSourceOption {
	return func(t *tokenSource) {
		t.noWrap = true
	}
}
