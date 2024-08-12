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
	accessToken string
	tokenType   string
	expiry      time.Time
	remoteToken bool // flag to indicate if the current token was retrieved from a remote source

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
	if !t.Valid() {
		if !t.cloudEnabled {
			return nil, ErrTokenExpired
		}
		if err := t.retrieveToken(ctx); err != nil {
			return nil, err
		}
	}

	return &oauth2.Token{
		AccessToken: t.accessToken,
		TokenType:   t.tokenType,
		Expiry:      t.expiry,
	}, nil
}

// Valid reports whether t is non-nil, has an AccessToken, and is not expired.
// If cloud env is set, it also checks if the token was retrieved from a remote source.
// This way we always retrieve token from remote source if cloud env is set
// provided the current token is not acquired from a remote source.
func (t *Token) Valid() bool {
	if t.cloudEnabled && !t.remoteToken {
		return false
	}

	return t != nil && t.accessToken != "" && !t.expired()
}

// timeNow is time.Now but pulled out as a variable for tests.
var timeNow = time.Now

// expired reports whether the token is expired.
// t must be non-nil.
func (t *Token) expired() bool {
	if t.expiry.IsZero() {
		return false
	}

	return t.expiry.Round(0).Add(-defaultExpiryDelta).Before(timeNow())
}

func (t *Token) retrieveToken(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.Valid() {
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
	t.accessToken = oauthResp.AccessToken
	if oauthResp.Expires == nil {
		t.expiry = time.Time{}
	} else {
		t.expiry = *oauthResp.Expires
	}
	return nil
}
