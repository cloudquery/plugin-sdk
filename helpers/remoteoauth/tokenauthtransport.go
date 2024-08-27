package remoteoauth

import (
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

// TokenAuthTransport is a custom http.RoundTripper to inject the OAuth2 token
type TokenAuthTransport struct {
	TokenSource oauth2.TokenSource // required

	BaseTransport http.RoundTripper
	Rewriter      RewriterFunc
}

var errNilToken = errors.New("remoteoauth: nil token")

// RewriterFunc is a function that can be used to rewrite the request before it is sent
type RewriterFunc func(*http.Request, oauth2.Token)

// RoundTrip executes a single HTTP transaction and injects the token into the request header
func (t *TokenAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	reqCopy := req.Clone(req.Context())

	token, err := t.TokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("remoteoauth: failed to get token: %w", err)
	}
	if token == nil {
		return nil, errNilToken
	}

	if t.Rewriter != nil {
		t.Rewriter(reqCopy, *token)
	} else {
		// Inject the token as the Authorization header
		token.SetAuthHeader(reqCopy)
	}

	transport := t.BaseTransport
	if transport == nil {
		transport = http.DefaultTransport
	}

	return transport.RoundTrip(reqCopy)
}

// NewAuthTransport creates a new TokenAuthTransport with the provided TokenSource
func NewAuthTransport(ts oauth2.TokenSource) http.RoundTripper {
	return &TokenAuthTransport{
		TokenSource: ts,
	}
}
