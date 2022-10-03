package clients

import (
	"context"
	"os"
	"testing"

	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
)

var newSourceClientTestCases = []specs.Source{
	{
		Name:     "test",
		Registry: specs.RegistryGithub,
		Path:     "cloudquery/test",
		Version:  "v1.1.5",
	},
	{
		Name:     "test",
		Registry: specs.RegistryGithub,
		Path:     "yevgenypats/test",
		Version:  "v1.0.1",
	},
}

func TestSourceClient(t *testing.T) {
	ctx := context.Background()
	l := zerolog.New(zerolog.NewTestWriter(t)).Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel)
	for _, tc := range newSourceClientTestCases {
		t.Run(tc.Path+"_"+tc.Version, func(t *testing.T) {
			dirName := t.TempDir()
			c, err := NewSourceClient(ctx, tc.Registry, tc.Path, tc.Version, WithSourceLogger(l), WithSourceDirectory(dirName))
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := c.Terminate(); err != nil {
					t.Fatalf("failed to terminate source client: %v", err)
				}
			}()
			tables, err := c.GetTables(ctx)
			if err != nil {
				t.Fatal("failed to get tables", err)
			}
			if len(tables) != 1 {
				t.Fatal("expected 1 table got ", len(tables))
			}
		})
	}
}
