package remoteoauth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

type TokenOption func(*Token)

func WithAccessToken(token, tokenType string, expiry time.Time) TokenOption {
	return func(t *Token) {
		t.currentToken.AccessToken = token
		t.currentToken.TokenType = tokenType
		t.currentToken.Expiry = expiry
	}
}

func (t *Token) initCloudOpts() error {
	_, t.cloudEnabled = os.LookupEnv("CQ_CLOUD")
	if !t.cloudEnabled {
		return nil
	}

	t.apiToken = os.Getenv("CLOUDQUERY_API_KEY")
	if t.apiToken == "" {
		return errors.New("CLOUDQUERY_API_KEY missing")
	}
	t.apiURL = os.Getenv("CLOUDQUERY_API_URL")
	if t.apiURL == "" {
		t.apiURL = "https://api.cloudquery.io"
	}

	t.teamName = os.Getenv("_CQ_TEAM_NAME")
	if t.teamName == "" {
		return errors.New("_CQ_TEAM_NAME missing")
	}
	t.syncName = os.Getenv("_CQ_SYNC_NAME")
	syncRunID := os.Getenv("_CQ_SYNC_RUN_ID")
	testConnID := os.Getenv("_CQ_SYNC_TEST_CONNECTION_ID")
	if testConnID == "" && syncRunID == "" {
		return errors.New("_CQ_SYNC_TEST_CONNECTION_ID or _CQ_SYNC_RUN_ID missing")
	} else if testConnID != "" && syncRunID != "" {
		return errors.New("_CQ_SYNC_TEST_CONNECTION_ID and _CQ_SYNC_RUN_ID are mutually exclusive")
	}

	var err error
	if syncRunID != "" {
		if t.syncName == "" {
			return errors.New("_CQ_SYNC_NAME missing")
		}

		t.syncRunUUID, err = uuid.Parse(syncRunID)
		if err != nil {
			return fmt.Errorf("_CQ_SYNC_RUN_ID is not a valid UUID: %w", err)
		}
	}
	if testConnID != "" {
		if t.syncName != "" {
			return errors.New("_CQ_SYNC_NAME should be empty")
		}

		t.testConnUUID, err = uuid.Parse(testConnID)
		if err != nil {
			return fmt.Errorf("_CQ_SYNC_TEST_CONNECTION_ID is not a valid UUID: %w", err)
		}
		t.isTestConnection = true
	}

	connectorID := os.Getenv("_CQ_CONNECTOR_ID")
	if connectorID == "" {
		return errors.New("_CQ_CONNECTOR_ID missing")
	}
	t.connectorUUID, err = uuid.Parse(connectorID)
	if err != nil {
		return fmt.Errorf("_CQ_CONNECTOR_ID is not a valid UUID: %w", err)
	}
	return nil
}
