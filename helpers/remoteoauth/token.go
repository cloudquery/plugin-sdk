package remoteoauth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	cloudquery_api "github.com/cloudquery/cloudquery-api-go"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type Token struct {
	AccessToken string
	TokenType   string
	Expiry      time.Time
	remoteToken bool // was the current token retrieved from a remote source

	apiClient *cloudquery_api.ClientWithResponses
	mu        sync.Mutex

	cloudEnabled     bool
	apiURL           string
	apiToken         string
	teamName         string
	syncName         string
	testConnUUID     uuid.UUID
	syncRunUUID      uuid.UUID
	connectorUUID    uuid.UUID
	isTestConnection bool
}

var (
	ErrTokenExpired = errors.New("token expired and cloud env is not set")

	_ oauth2.TokenSource = (*Token)(nil)
)

// defaultExpiryDelta determines how earlier a token should be considered
// expired than its actual expiration time. It is used to avoid late
// expirations due to client-server time mismatches.
const defaultExpiryDelta = 10 * time.Second

func NewToken(opts ...TokenOption) (*Token, error) {
	t := &Token{}
	for _, opt := range opts {
		opt(t)
	}
	err := t.initCloudOpts()
	if err != nil {
		return nil, err
	}
	if t.cloudEnabled {
		t.apiClient, err = cloudquery_api.NewClientWithResponses(t.apiURL,
			cloudquery_api.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
				req.Header.Set("Authorization", "Bearer "+t.apiToken)
				return nil
			}))
		if err != nil {
			return nil, fmt.Errorf("failed to create api client: %w", err)
		}
	}

	return t, nil
}

// Token returns the cached token if not expired, or a new token from the remote source.
func (t *Token) Token() (*oauth2.Token, error) {
	return t.TokenWithContext(context.TODO())
}

// TokenWithContext returns the cached token if not expired, or a new token from the remote source
// using the given context.
func (t *Token) TokenWithContext(ctx context.Context) (*oauth2.Token, error) {
	if t.cloudEnabled && !t.remoteToken {
		// Always retrieve token from remote source if cloud env is set
		// and the current token is not acquired from a remote source.
		if err := t.retrieveToken(ctx); err != nil {
			return nil, err
		}
	} else if !t.Valid() {
		if !t.cloudEnabled {
			return nil, ErrTokenExpired
		}
		if err := t.retrieveToken(ctx); err != nil {
			return nil, err
		}
	}

	return &oauth2.Token{
		AccessToken: t.AccessToken,
		TokenType:   t.TokenType,
		Expiry:      t.Expiry,
	}, nil
}

// Valid reports whether t is non-nil, has an AccessToken, and is not expired.
func (t *Token) Valid() bool {
	return t != nil && t.AccessToken != "" && !t.expired()
}

// timeNow is time.Now but pulled out as a variable for tests.
var timeNow = time.Now

// expired reports whether the token is expired.
// t must be non-nil.
func (t *Token) expired() bool {
	if t.Expiry.IsZero() {
		return false
	}

	return t.Expiry.Round(0).Add(-defaultExpiryDelta).Before(timeNow())
}

func (t *Token) retrieveToken(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.remoteToken && t.Valid() {
		// if another goroutine has updated the token, return
		return nil
	}

	var oauthResp *cloudquery_api.ConnectorCredentialsResponseOAuth
	if !t.isTestConnection {
		resp, err := t.apiClient.GetSyncRunConnectorCredentialsWithResponse(ctx, t.teamName, t.syncName, t.syncRunUUID, t.connectorUUID)
		if err != nil {
			return fmt.Errorf("failed to get sync run connector credentials: %w", err)
		}
		if resp.StatusCode() != http.StatusOK {
			if resp.JSON422 != nil {
				return fmt.Errorf("failed to get sync run connector credentials: %s", resp.JSON422.Message)
			}
			return fmt.Errorf("failed to get sync run connector credentials: %s", resp.Status())
		}
		oauthResp = resp.JSON200.Oauth
	} else {
		resp, err := t.apiClient.GetTestConnectionConnectorCredentialsWithResponse(ctx, t.teamName, t.testConnUUID, t.connectorUUID)
		if err != nil {
			return fmt.Errorf("failed to get test connection connector credentials: %w", err)
		}
		if resp.StatusCode() != http.StatusOK {
			if resp.JSON422 != nil {
				return fmt.Errorf("failed to get test connection connector credentials: %s", resp.JSON422.Message)
			}
			return fmt.Errorf("failed to get test connection connector credentials: %s", resp.Status())
		}
		oauthResp = resp.JSON200.Oauth
	}

	if oauthResp == nil {
		return fmt.Errorf("missing oauth credentials in response")
	}
	t.remoteToken = true
	t.AccessToken = oauthResp.AccessToken
	if oauthResp.Expires == nil {
		t.Expiry = time.Time{}
	} else {
		t.Expiry = *oauthResp.Expires
	}
	return nil
}
