package remoteoauth

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

// TokenAuthEditor returns a custom RequestEditorFn to inject the OAuth2 token
func TokenAuthEditor(tokenSource oauth2.TokenSource) func(context.Context, *http.Request) error {
	return func(_ context.Context, req *http.Request) error {
		token, err := tokenSource.Token()
		if err != nil {
			return fmt.Errorf("remoteoauth: failed to get token: %w", err)
		}
		if token == nil {
			return errNilToken
		}

		// Inject the token as the Authorization header
		token.SetAuthHeader(req)
		return nil
	}
}
