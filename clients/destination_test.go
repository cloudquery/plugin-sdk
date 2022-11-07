package clients

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/cloudquery/plugin-sdk/v1/specs"
	"github.com/rs/zerolog"
)

var newDestinationClientTestCases = []specs.Source{
	{
		Name:     "test",
		Registry: specs.RegistryGithub,
		Path:     "cloudquery/test",
		Version:  "v1.1.0",
	},
	{
		Name:     "test",
		Registry: specs.RegistryGithub,
		Path:     "yevgenypats/test",
		Version:  "v1.0.1",
	},
}

// TestDestinationClient mostly checks the download and spawn logic. it doesn't call all methods as those are
// tested under serve/tests
func TestDestinationClient(t *testing.T) {
	ctx := context.Background()
	l := zerolog.New(zerolog.NewTestWriter(t)).Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel)
	for _, tc := range newDestinationClientTestCases {
		t.Run(tc.Path+"_"+tc.Version, func(t *testing.T) {
			dirName := t.TempDir()
			c, err := NewDestinationClient(ctx, tc.Registry, tc.Path, tc.Version, WithDestinationLogger(l), WithDestinationDirectory(dirName))
			if err != nil {
				if strings.HasPrefix(err.Error(), "destination plugin protocol version") {
					// this also means success as in this tests we just want to make sure we were able to download and spawn the plugin
					return
				}
				t.Fatal(err)
			}
			defer func() {
				if err := c.Terminate(); err != nil {
					t.Logf("failed to terminate destination client: %v", err)
				}
			}()
			if err := c.Initialize(ctx, specs.Destination{}); err != nil {
				t.Fatal(err)
			}
			name, err := c.Name(ctx)
			if err != nil {
				t.Fatal("failed to get name", err)
			}
			if name != "test" {
				t.Fatal("expected name to be test got ", name)
			}
		})
	}
}
