package remoteoauth

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestInitCloudOpts(t *testing.T) {
	validUUID := uuid.NewString()

	cases := []struct {
		name           string
		envs           map[string]string
		expectError    bool
		expectCloud    bool
		expectTestConn bool
	}{
		{
			name: "no envs",
		},
		{
			name: "cloud env",
			envs: map[string]string{
				"CQ_CLOUD":           "1",
				"CLOUDQUERY_API_KEY": "the-key",
				"_CQ_TEAM_NAME":      "the-team",
				"_CQ_SYNC_NAME":      "the-sync",
				"_CQ_SYNC_RUN_ID":    validUUID,
				"_CQ_CONNECTOR_ID":   validUUID,
			},
			expectCloud: true,
		},
		{
			name: "cloud env test conn",
			envs: map[string]string{
				"CQ_CLOUD":                    "1",
				"CLOUDQUERY_API_KEY":          "the-key",
				"_CQ_TEAM_NAME":               "the-team",
				"_CQ_SYNC_TEST_CONNECTION_ID": validUUID,
				"_CQ_CONNECTOR_ID":            validUUID,
			},
			expectCloud:    true,
			expectTestConn: true,
		},
		{
			name: "missing cq_cloud with everything set",
			envs: map[string]string{
				"CLOUDQUERY_API_KEY": "the-key",
				"_CQ_TEAM_NAME":      "the-team",
				"_CQ_SYNC_NAME":      "the-sync",
				"_CQ_SYNC_RUN_ID":    validUUID,
				"_CQ_CONNECTOR_ID":   validUUID,
			},
			expectCloud: false,
		},
		{
			name: "missing cq_cloud with missing api key",
			envs: map[string]string{
				"_CQ_TEAM_NAME":    "the-team",
				"_CQ_SYNC_NAME":    "the-sync",
				"_CQ_SYNC_RUN_ID":  validUUID,
				"_CQ_CONNECTOR_ID": validUUID,
			},
			expectCloud: false,
		},
		{
			name: "missing cq_cloud with missing sync name",
			envs: map[string]string{
				"CLOUDQUERY_API_KEY": "the-key",
				"_CQ_TEAM_NAME":      "the-team",
				"_CQ_SYNC_NAME":      "the-sync",
				"_CQ_SYNC_RUN_ID":    validUUID,
				"_CQ_CONNECTOR_ID":   validUUID,
			},
			expectCloud: false,
		},
		{
			name: "cloud env missing api key",
			envs: map[string]string{
				"CQ_CLOUD":                    "1",
				"_CQ_TEAM_NAME":               "the-team",
				"_CQ_SYNC_TEST_CONNECTION_ID": validUUID,
				"_CQ_CONNECTOR_ID":            validUUID,
			},
			expectError: true,
		},
		{
			name: "cloud env missing sync name",
			envs: map[string]string{
				"CQ_CLOUD":           "1",
				"CLOUDQUERY_API_KEY": "the-key",
				"_CQ_TEAM_NAME":      "the-team",
				"_CQ_SYNC_RUN_ID":    validUUID,
				"_CQ_CONNECTOR_ID":   validUUID,
			},
			expectError: true,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)
			for k, v := range tc.envs {
				t.Setenv(k, v)
			}
			tok, err := NewTokenSource(withNoWrap())
			if tc.expectError {
				r.Error(err)
				return
			}
			r.NoError(err)
			if tc.expectCloud {
				ts := tok.(*cloudTokenSource)
				r.Equal(tc.expectTestConn, ts.isTestConnection)
				return
			}
			rt := reflect.TypeOf(tok)
			r.Equal("oauth2.staticTokenSource", rt.String())
		})
	}
}
