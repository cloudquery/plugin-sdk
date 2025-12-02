package remoteoauth

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

func TestLocalTokenAccess(t *testing.T) {
	r := require.New(t)
	_, cloud := os.LookupEnv("CQ_CLOUD")
	r.False(cloud, "CQ_CLOUD should not be set")
	tok, err := NewTokenSource(WithToken(oauth2.Token{AccessToken: "token", TokenType: "bearer"}))
	r.NoError(err)
	tk, err := tok.Token()
	r.NoError(err)
	r.True(tk.Valid())
	r.Equal("token", tk.AccessToken)
}

func TestLocalTokenAccessWithDeprecatedTokenOpt(t *testing.T) {
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
