package clients

import (
	"context"
	"encoding/json"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
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

func TestDestinationClientWriteReturnsCorrectError(t *testing.T) {
	ctx := context.Background()
	l := zerolog.New(zerolog.NewTestWriter(t)).Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel)
	dirName := t.TempDir()
	c, err := NewDestinationClient(ctx, specs.RegistryGithub, "cloudquery/sqlite", "v1.0.11", WithDestinationLogger(l), WithDestinationDirectory(dirName))
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := c.Terminate(); err != nil {
			t.Logf("failed to terminate destination client: %v", err)
		}
	}()
	sqliteSpec := struct {
		connectionString string
	}{connectionString: path.Join(dirName, "test.sql")}
	if err := c.Initialize(ctx, specs.Destination{Spec: sqliteSpec}); err != nil {
		t.Fatal(err)
	}

	_, err = c.Name(ctx)
	if err != nil {
		t.Fatal("failed to get name", err)
	}

	columns := []schema.Column{{Name: "int", Type: schema.TypeInt}}
	tables := schema.Tables{&schema.Table{Name: "test-1", Columns: columns}, &schema.Table{Name: "test-2", Columns: columns}}
	resource1 := schema.Resource{Item: map[string]any{"int": 1}, Table: tables[0]}
	destResource1, _ := json.Marshal(resource1.ToDestinationResource())
	resource2 := schema.Resource{Item: map[string]any{"int": 1}, Table: tables[1]}
	destResource2, _ := json.Marshal(resource2.ToDestinationResource())
	resourcesChannel := make(chan []byte)
	go func() {
		defer close(resourcesChannel)
		// we need to stream enough data to the server so it at least starts processing it and return the relevant error
		for i := 1; i < 100000; i++ {
			resourcesChannel <- destResource1
			resourcesChannel <- destResource2
			resourcesChannel <- destResource1
			resourcesChannel <- destResource2
		}
	}()
	sourceSpec := specs.Source{
		Name: "TestDestinationClientWriteReturnsCorrectError",
	}
	err = c.Write2(ctx, sourceSpec, tables, time.Now().UTC(), resourcesChannel)
	require.ErrorContains(t, err, "context canceled")
}
