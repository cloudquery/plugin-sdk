package clients

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/cloudquery/plugin-sdk/specs"
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

func TestDestinationClient(t *testing.T) {
	ctx := context.Background()
	l := zerolog.New(zerolog.NewTestWriter(t)).Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel)
	for _, tc := range newDestinationClientTestCases {
		t.Run(tc.Path+"_"+tc.Version, func(t *testing.T) {
			dirName := t.TempDir()
			localPath := path.Join(dirName, "plugin")
			if err := DownloadPluginFromGithub(ctx, localPath, tc.Path, tc.Version, PluginTypeDestination); err != nil {
				t.Fatal(err)
			}
			c, err := NewManagedDestinationClient(ctx, localPath, WithDestinationLogger(l))
			if err != nil {
				t.Fatal(err)
			}
			defer c.Close()
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
