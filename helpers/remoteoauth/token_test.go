package remoteoauth

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const testAPIKey = "test-key"

func TestLocalTokenAccess(t *testing.T) {
	r := require.New(t)
	_, cloud := os.LookupEnv("CQ_CLOUD")
	r.False(cloud, "CQ_CLOUD should not be set")
	tok, err := NewTokenSource(WithAccessToken("token", "bearer", time.Time{}))
	r.NoError(err)
	tk, err := tok.Token()
	r.NoError(err)
	r.True(tk.Valid())
	r.Equal("token", tk.AccessToken)
}

func TestFirstLocalTokenAccess(t *testing.T) {
	runID := uuid.NewString()
	connID := uuid.NewString()
	testURL := setupMockTokenServer(t, map[string]string{
		"/teams/the-team/syncs/the-sync/runs/" + runID + "/connector/" + connID + "/credentials": `{"oauth":{"access_token":"new-token"}}`,
	})
	setEnvs(t, map[string]string{
		"CQ_CLOUD":           "1",
		"CLOUDQUERY_API_URL": testURL,
		"CLOUDQUERY_API_KEY": testAPIKey,
		"_CQ_TEAM_NAME":      "the-team",
		"_CQ_SYNC_NAME":      "the-sync",
		"_CQ_SYNC_RUN_ID":    runID,
		"_CQ_CONNECTOR_ID":   connID,
	})
	r := require.New(t)
	tok, err := NewTokenSource(WithAccessToken("token", "bearer", time.Time{}))
	r.NoError(err)
	tk, err := tok.Token()
	r.NoError(err)
	r.True(tk.Valid())
	r.Equal("new-token", tk.AccessToken)
}

func TestInvalidAPIKeyTokenAccess(t *testing.T) {
	runID := uuid.NewString()
	connID := uuid.NewString()
	testURL := setupMockTokenServer(t, nil)
	setEnvs(t, map[string]string{
		"CQ_CLOUD":           "1",
		"CLOUDQUERY_API_URL": testURL,
		"CLOUDQUERY_API_KEY": "invalid",
		"_CQ_TEAM_NAME":      "the-team",
		"_CQ_SYNC_NAME":      "the-sync",
		"_CQ_SYNC_RUN_ID":    runID,
		"_CQ_CONNECTOR_ID":   connID,
	})
	r := require.New(t)
	tok, err := NewTokenSource(WithAccessToken("token", "bearer", time.Time{}))
	r.NoError(err)
	tk, err := tok.Token()
	r.Nil(tk)
	r.False(tk.Valid())
	r.ErrorContains(err, "failed to get sync run connector credentials")
}

func TestSyncRunTokenAccess(t *testing.T) {
	runID := uuid.NewString()
	connID := uuid.NewString()
	testURL := setupMockTokenServer(t, map[string]string{
		"/teams/the-team/syncs/the-sync/runs/" + runID + "/connector/" + connID + "/credentials": `{"oauth":{"access_token":"new-token"}}`,
	})
	setEnvs(t, map[string]string{
		"CQ_CLOUD":           "1",
		"CLOUDQUERY_API_URL": testURL,
		"CLOUDQUERY_API_KEY": testAPIKey,
		"_CQ_TEAM_NAME":      "the-team",
		"_CQ_SYNC_NAME":      "the-sync",
		"_CQ_SYNC_RUN_ID":    runID,
		"_CQ_CONNECTOR_ID":   connID,
	})
	r := require.New(t)
	tok, err := NewTokenSource()
	r.NoError(err)
	tk, err := tok.Token()
	r.NoError(err)
	r.True(tk.Valid())
	r.Equal("new-token", tk.AccessToken)
}

func TestTestConnectionTokenAccess(t *testing.T) {
	testID := uuid.NewString()
	connID := uuid.NewString()
	testURL := setupMockTokenServer(t, map[string]string{
		"/teams/the-team/syncs/test-connections/" + testID + "/connector/" + connID + "/credentials": `{"oauth":{"access_token":"new-token"}}`,
	})
	setEnvs(t, map[string]string{
		"CQ_CLOUD":                    "1",
		"CLOUDQUERY_API_URL":          testURL,
		"CLOUDQUERY_API_KEY":          testAPIKey,
		"_CQ_TEAM_NAME":               "the-team",
		"_CQ_SYNC_TEST_CONNECTION_ID": testID,
		"_CQ_CONNECTOR_ID":            connID,
	})
	r := require.New(t)
	tok, err := NewTokenSource(WithAccessToken("token", "bearer", time.Time{}))
	r.NoError(err)
	tk, err := tok.Token()
	r.NoError(err)
	r.True(tk.Valid())
	r.Equal("new-token", tk.AccessToken)
}

func setEnvs(t *testing.T, envs map[string]string) {
	t.Helper()
	for k, v := range envs {
		t.Setenv(k, v)
	}
}

func setupMockTokenServer(t *testing.T, responses map[string]string) string {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a := r.Header.Get("Authorization"); a != "Bearer "+testAPIKey {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		resp, ok := responses[r.URL.Path]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(resp))
	}))
	t.Cleanup(func() {
		ts.Close()
	})
	return ts.URL
}
