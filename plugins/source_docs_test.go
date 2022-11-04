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
			{
				Name:            "id_col",
				Type:            schema.TypeInt,
				CreationOptions: schema.ColumnCreationOptions{PrimaryKey: true},
			},
			{
				Name:            "id_col2",
				Type:            schema.TypeInt,
				CreationOptions: schema.ColumnCreationOptions{PrimaryKey: true},
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
			{
				Name:        "relation_table2",
				Description: "Description for second relational table",
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

func TestGenerateSourcePluginDocs(t *testing.T) {
	p := NewSourcePlugin("test", "v1.0.0", testTables, newTestExecutionClient)

	t.Run("Markdown", func(t *testing.T) {
		tmpdir := t.TempDir()

		err := p.GenerateSourcePluginDocs(tmpdir, SourceDocsFormatMarkdown)
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
	})

	t.Run("JSON", func(t *testing.T) {
		tmpdir := t.TempDir()

		err := p.GenerateSourcePluginDocs(tmpdir, SourceDocsFormatJSON)
		if err != nil {
			t.Fatalf("unexpected error calling GenerateSourcePluginDocs: %v", err)
		}

		expectFiles := []string{"__tables.json"}
		for _, exp := range expectFiles {
			t.Run(exp, func(t *testing.T) {
				output := path.Join(tmpdir, exp)
				got, err := os.ReadFile(output)
				require.NoError(t, err)
				cupaloy.SnapshotT(t, got)
			})
		}
	})
}
