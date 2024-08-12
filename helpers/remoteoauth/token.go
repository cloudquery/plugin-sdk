package remoteoauth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"

	cloudquery_api "github.com/cloudquery/cloudquery-api-go"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

func NewTokenSource(opts ...TokenSourceOption) (oauth2.TokenSource, error) {
	t := &tokenSource{}
	for _, opt := range opts {
		opt(t)
	}

	if _, cloudEnabled := os.LookupEnv("CQ_CLOUD"); !cloudEnabled {
		return oauth2.StaticTokenSource(&t.currentToken), nil
	}

	cloudToken, err := newCloudTokenSource()
	if err != nil {
		return nil, err
	}
	if t.noWrap {
		return cloudToken, nil
	}

	return oauth2.ReuseTokenSource(nil, cloudToken), nil
}

type tokenSource struct {
	currentToken oauth2.Token
	noWrap       bool
}

type cloudTokenSource struct {
	apiClient *cloudquery_api.ClientWithResponses
	mu        sync.Mutex // protects against multiple refresh calls

	apiURL           string
	apiToken         string
	teamName         string
	syncName         string
	testConnUUID     uuid.UUID
	syncRunUUID      uuid.UUID
	connectorUUID    uuid.UUID
	isTestConnection bool
}

var _ oauth2.TokenSource = (*cloudTokenSource)(nil)

func newCloudTokenSource() (oauth2.TokenSource, error) {
	t := &cloudTokenSource{}

	err := t.initCloudOpts()
	if err != nil {
		return nil, err
	}

	t.apiClient, err = cloudquery_api.NewClientWithResponses(t.apiURL,
		cloudquery_api.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+t.apiToken)
			return nil
		}))
	if err != nil {
		return nil, fmt.Errorf("failed to create api client: %w", err)
	}

	return t, nil
}

// Token returns the cached token if not expired, or a new token from the remote source.
func (t *cloudTokenSource) Token() (*oauth2.Token, error) {
	return t.TokenWithContext(context.TODO())
}

// TokenWithContext returns a new token from the remote source using the given context.
func (t *cloudTokenSource) TokenWithContext(ctx context.Context) (*oauth2.Token, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	var oauthResp *cloudquery_api.ConnectorCredentialsResponseOAuth
	if !t.isTestConnection {
		resp, err := t.apiClient.GetSyncRunConnectorCredentialsWithResponse(ctx, t.teamName, t.syncName, t.syncRunUUID, t.connectorUUID)
		if err != nil {
			return nil, fmt.Errorf("failed to get sync run connector credentials: %w", err)
		}
		if resp.StatusCode() != http.StatusOK {
			if resp.JSON422 != nil {
				return nil, fmt.Errorf("failed to get sync run connector credentials: %s", resp.JSON422.Message)
			}
			return nil, fmt.Errorf("failed to get sync run connector credentials: %s", resp.Status())
		}
		oauthResp = resp.JSON200.Oauth
	} else {
		resp, err := t.apiClient.GetTestConnectionConnectorCredentialsWithResponse(ctx, t.teamName, t.testConnUUID, t.connectorUUID)
		if err != nil {
			return nil, fmt.Errorf("failed to get test connection connector credentials: %w", err)
		}
		if resp.StatusCode() != http.StatusOK {
			if resp.JSON422 != nil {
				return nil, fmt.Errorf("failed to get test connection connector credentials: %s", resp.JSON422.Message)
			}
			return nil, fmt.Errorf("failed to get test connection connector credentials: %s", resp.Status())
		}
		oauthResp = resp.JSON200.Oauth
	}

	if oauthResp == nil {
		return nil, fmt.Errorf("missing oauth credentials in response")
	}

	tok := &oauth2.Token{
		AccessToken: oauthResp.AccessToken,
	}
	if oauthResp.Expires != nil {
		tok.Expiry = *oauthResp.Expires
	}
	return tok, nil
}

func (t *cloudTokenSource) initCloudOpts() error {
	var allErr error

	t.apiToken = os.Getenv("CLOUDQUERY_API_KEY")
	if t.apiToken == "" {
		allErr = errors.Join(allErr, errors.New("CLOUDQUERY_API_KEY missing"))
	}
	t.apiURL = os.Getenv("CLOUDQUERY_API_URL")
	if t.apiURL == "" {
		t.apiURL = "https://api.cloudquery.io"
	}

	t.teamName = os.Getenv("_CQ_TEAM_NAME")
	if t.teamName == "" {
		allErr = errors.Join(allErr, errors.New("_CQ_TEAM_NAME missing"))
	}
	t.syncName = os.Getenv("_CQ_SYNC_NAME")
	syncRunID := os.Getenv("_CQ_SYNC_RUN_ID")
	testConnID := os.Getenv("_CQ_SYNC_TEST_CONNECTION_ID")
	if testConnID == "" && syncRunID == "" {
		allErr = errors.Join(allErr, errors.New("_CQ_SYNC_TEST_CONNECTION_ID or _CQ_SYNC_RUN_ID missing"))
	} else if testConnID != "" && syncRunID != "" {
		allErr = errors.Join(allErr, errors.New("_CQ_SYNC_TEST_CONNECTION_ID and _CQ_SYNC_RUN_ID are mutually exclusive"))
	}

	var err error
	if syncRunID != "" {
		if t.syncName == "" {
			allErr = errors.Join(allErr, errors.New("_CQ_SYNC_NAME missing"))
		}

		t.syncRunUUID, err = uuid.Parse(syncRunID)
		if err != nil {
			allErr = errors.Join(allErr, fmt.Errorf("_CQ_SYNC_RUN_ID is not a valid UUID: %w", err))
		}
	}
	if testConnID != "" {
		if t.syncName != "" {
			allErr = errors.Join(allErr, errors.New("_CQ_SYNC_NAME should be empty"))
		}

		t.testConnUUID, err = uuid.Parse(testConnID)
		if err != nil {
			allErr = errors.Join(allErr, fmt.Errorf("_CQ_SYNC_TEST_CONNECTION_ID is not a valid UUID: %w", err))
		}
		t.isTestConnection = true
	}

	connectorID := os.Getenv("_CQ_CONNECTOR_ID")
	if connectorID == "" {
		allErr = errors.Join(allErr, errors.New("_CQ_CONNECTOR_ID missing"))
	} else {
		t.connectorUUID, err = uuid.Parse(connectorID)
		if err != nil {
			allErr = errors.Join(allErr, fmt.Errorf("_CQ_CONNECTOR_ID is not a valid UUID: %w", err))
		}
	}
	return allErr
}
