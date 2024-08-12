package remoteoauth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type TokenOption func(*Token)

func WithAccessToken(token, tokenType string, expiry time.Time) TokenOption {
	return func(t *Token) {
		t.currentToken = oauth2.Token{
			AccessToken: token,
			TokenType:   tokenType,
			Expiry:      expiry,
		}
	}
}

func (t *Token) initCloudOpts() error {
	_, t.cloudEnabled = os.LookupEnv("CQ_CLOUD")
	if !t.cloudEnabled {
		return nil
	}

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
