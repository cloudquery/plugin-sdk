package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

type RemoteOAuthTokenOption func(*RemoteOAuthToken)

func WithAccessToken(token, tokenType string, expiry time.Time) RemoteOAuthTokenOption {
	return func(t *RemoteOAuthToken) {
		t.AccessToken = token
		t.TokenType = tokenType
		t.Expiry = expiry
	}
}

func WithCloudEnv() RemoteOAuthTokenOption {
	return func(t *RemoteOAuthToken) {
		_, t.cloudEnabled = os.LookupEnv("CQ_CLOUD")
		if !t.cloudEnabled {
			return
		}

		t.teamName = os.Getenv("_CQ_TEAM_NAME")
		t.syncName = os.Getenv("_CQ_SYNC_NAME")
		t.syncRunID = os.Getenv("_CQ_SYNC_RUN_ID")
		t.testConnID = os.Getenv("_CQ_SYNC_TEST_CONNECTION_ID")
		t.connectorID = os.Getenv("_CQ_CONNECTOR_ID")
		t.apiToken = os.Getenv("CLOUDQUERY_API_KEY")
		t.apiURL = os.Getenv("CLOUDQUERY_API_URL")
		if t.apiURL == "" {
			t.apiURL = "https://api.cloudquery.io"
		}
	}
}

func (t *RemoteOAuthToken) validateCloudOpts() error {
	var err error
	if t.apiURL == "" {
		return errors.New("CLOUDQUERY_API_URL is empty")
	}
	if t.apiToken == "" {
		return errors.New("CLOUDQUERY_API_KEY missing")
	}
	if t.teamName == "" {
		return errors.New("_CQ_TEAM_NAME missing")
	}
	if t.testConnID == "" && t.syncRunID == "" {
		return errors.New("_CQ_SYNC_TEST_CONNECTION_ID or CQ_SYNC_RUN_ID missing")
	} else if t.testConnID != "" && t.syncRunID != "" {
		return errors.New("_CQ_SYNC_TEST_CONNECTION_ID and CQ_SYNC_RUN_ID are mutually exclusive")
	}
	if t.syncRunID != "" {
		if t.syncName == "" {
			return errors.New("_CQ_SYNC_NAME missing")
		}

		t.syncRunUUID, err = uuid.Parse(t.syncRunID)
		if err != nil {
			return fmt.Errorf("_CQ_SYNC_RUN_ID is not a valid UUID: %w", err)
		}
	}
	if t.testConnID != "" {
		if t.syncName != "" {
			return errors.New("_CQ_SYNC_NAME should be empty")
		}

		t.testConnUUID, err = uuid.Parse(t.testConnID)
		if err != nil {
			return fmt.Errorf("_CQ_SYNC_TEST_CONNECTION_ID is not a valid UUID: %w", err)
		}
	}

	if t.connectorID == "" {
		return errors.New("_CQ_CONNECTOR_ID missing")
	}
	t.connectorUUID, err = uuid.Parse(t.connectorID)
	if err != nil {
		return fmt.Errorf("_CQ_CONNECTOR_ID is not a valid UUID: %w", err)
	}
	return nil
}
