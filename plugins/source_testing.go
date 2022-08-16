package plugins

import (
	"context"
	"testing"

	"github.com/cloudquery/faker/v3"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/georgysavva/scany/pgxscan"
)

type ResourceTestCase struct {
	Plugin *SourcePlugin
	Spec   specs.Source
	// ParallelFetchingLimit limits parallel resources fetch at a time
	ParallelFetchingLimit uint64
	// SkipIgnoreInTest flag which detects if schema.Table or schema.Column should be ignored
	SkipIgnoreInTest bool
	// Verifiers are map from resource name to its verifiers.
	// If no verifiers specified for resource (resource name is not in key set of map),
	// non emptiness check of all columns in table and its relations will be performed.
	Verifiers map[string][]Verifier
}

// Verifier verifies tables specified by table schema (main table and its relations).
type Verifier func(t *testing.T, table *schema.Table, conn pgxscan.Querier, shouldSkipIgnoreInTest bool)

func init() {
	_ = faker.SetRandomMapAndSliceMinSize(1)
	_ = faker.SetRandomMapAndSliceMaxSize(1)
}

// type

func TestResource(t *testing.T, tc ResourceTestCase) {
	t.Parallel()
	t.Helper()

	// No need for configuration or db connection, get it out of the way first
	// testTableIdentifiersForProvider(t, resource.Provider)

	// l := testlog.New(t)
	// l.SetLevel(hclog.Info)
	// resource.Plugin.Logger = l
	resources := make(chan *schema.Resource)
	var fetchErr error

	// tc.Plugin.Logger = zerolog.New(zerolog.NewTestWriter(t))
	go func() {
		defer close(resources)
		fetchErr = tc.Plugin.Sync(context.Background(), tc.Spec, resources)
	}()

	for resource := range resources {
		validateResource(t, resource)
	}

	if fetchErr != nil {
		t.Fatal(fetchErr)
	}
}

func validateResource(t *testing.T, resource *schema.Resource) {
	t.Helper()
	for _, columnName := range resource.Table.Columns.Names() {
		if resource.Get(columnName) == nil && !resource.Table.Columns.Get(columnName).IgnoreInTests {
			t.Errorf("table: %s with unset column %s", resource.Table.Name, columnName)
		}
	}
}
