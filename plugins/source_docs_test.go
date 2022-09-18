package plugins

import (
	"os"
	"path"
	"testing"

	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/stretchr/testify/require"
)

var testTables = []*schema.Table{
	{
		Name:        "test_table",
		Description: "Description for test table",
		Columns: []schema.Column{
			{
				Name: "int_col",
				Type: schema.TypeInt,
			},
		},
		Relations: []*schema.Table{
			{
				Name:        "relation_table",
				Description: "Description for relational table",
				Columns: []schema.Column{
					{
						Name: "string_col",
						Type: schema.TypeString,
					},
				},
			},
		},
	},
}

// func (c *testExecutionClient) Logger() *zerolog.Logger {
// 	return &c.logger
// }

func TestGenerateSourcePluginDocs(t *testing.T) {
	tmpdir, tmpErr := os.MkdirTemp("", "docs_test_*")
	if tmpErr != nil {
		t.Fatalf("failed to create temporary directory: %v", tmpErr)
	}
	defer os.RemoveAll(tmpdir)

	p := NewSourcePlugin("test", "v1.0.0", testTables, newTestExecutionClient)
	err := p.GenerateSourcePluginDocs(tmpdir)
	if err != nil {
		t.Fatalf("unexpected error calling GenerateSourcePluginDocs: %v", err)
	}

	expectFiles := []string{"test_table.md", "relation_table.md", "README.md"}
	for _, exp := range expectFiles {
		t.Run(exp, func(t *testing.T) {
			output := path.Join(tmpdir, exp)
			got, err := os.ReadFile(output)
			require.NoError(t, err)
			cupaloy.SnapshotT(t, got)
		})
	}
}
