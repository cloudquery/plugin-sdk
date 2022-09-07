package docs

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/cloudquery/plugin-sdk/plugins"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/google/go-cmp/cmp"
	"github.com/rs/zerolog"
)

type testExecutionClient struct {
	logger zerolog.Logger
}

var testTables = []*schema.Table{
	{
		Name:        "test_table",
		Description: "Description for test table",
		Columns: []schema.Column{
			{
				Name:        "int_col",
				Type:        schema.TypeInt,
				Description: "Int column",
			},
		},
		Relations: []*schema.Table{
			{
				Name:        "relation_table",
				Description: "Description for relational table",
				Columns: []schema.Column{
					{
						Name:        "string_col",
						Type:        schema.TypeString,
						Description: "String column",
					},
				},
			},
		},
	},
}

var expectFiles = []struct {
	Name    string
	Content string
}{
	{
		Name: "test_table.md",
		Content: `
# Table: test_table
Description for test table
## Columns
| Name        | Type           | Description  |
| ------------- | ------------- | -----  |
|int_col|Int|Int column|
|_cq_id|UUID|Internal CQ ID of the row|
|_cq_fetch_time|Timestamp|Internal CQ row of when fetch was started (this will be the same for all rows in a single fetch)|
`,
	},
	{
		Name: "relation_table.md",
		Content: `
# Table: relation_table
Description for relational table
## Columns
| Name        | Type           | Description  |
| ------------- | ------------- | -----  |
|string_col|String|String column|
|_cq_id|UUID|Internal CQ ID of the row|
|_cq_fetch_time|Timestamp|Internal CQ row of when fetch was started (this will be the same for all rows in a single fetch)|
`,
	},
}

func (c *testExecutionClient) Logger() *zerolog.Logger {
	return &c.logger
}

func newTestExecutionClient(context.Context, zerolog.Logger, specs.Source) (schema.ClientMeta, error) {
	return &testExecutionClient{}, nil
}

func TestGenerateSourcePluginDocs(t *testing.T) {
	tmpdir, tmpErr := os.MkdirTemp("", "docs_test_*")
	if tmpErr != nil {
		t.Fatalf("failed to create temporary directory: %v", tmpErr)
	}
	defer os.RemoveAll(tmpdir)

	p := plugins.NewSourcePlugin("test", "v1.0.0", testTables, newTestExecutionClient)
	err := GenerateSourcePluginDocs(p, tmpdir)
	if err != nil {
		t.Fatalf("unexpected error calling GenerateSourcePluginDocs: %v", err)
	}

	for _, exp := range expectFiles {
		output := path.Join(tmpdir, exp.Name)
		got, err := ioutil.ReadFile(output)
		if err != nil {
			t.Fatalf("error reading %q: %v ", exp.Name, err)
		}

		if diff := cmp.Diff(string(got), exp.Content); diff != "" {
			t.Errorf("Generate docs for %q not as expected (+got, -want): %v", exp.Name, diff)
		}
	}
}
