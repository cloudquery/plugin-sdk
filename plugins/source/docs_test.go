//go:build !windows

package source

import (
	"os"
	"path"
	"testing"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/cloudquery/plugin-sdk/v3/schema"
	"github.com/stretchr/testify/require"
)

var testTables = []*schema.Table{
	{
		Name:        "test_table",
		Description: "Description for test table",
		Columns: []schema.Column{
			{
				Name: "int_col",
				Type: arrow.PrimitiveTypes.Int64,
			},
			{
				Name:       "id_col",
				Type:       arrow.PrimitiveTypes.Int64,
				PrimaryKey: true,
			},
			{
				Name:       "id_col2",
				Type:       arrow.PrimitiveTypes.Int64,
				PrimaryKey: true,
			},
		},
		Relations: []*schema.Table{
			{
				Name:        "relation_table",
				Description: "Description for relational table",
				Columns: []schema.Column{
					{
						Name: "string_col",
						Type: arrow.BinaryTypes.String,
					},
				},
				Relations: []*schema.Table{
					{
						Name:        "relation_relation_table_b",
						Description: "Description for relational table's relation",
						Columns: []schema.Column{
							{
								Name: "string_col",
								Type: arrow.BinaryTypes.String,
							},
						},
					},
					{
						Name:        "relation_relation_table_a",
						Description: "Description for relational table's relation",
						Columns: []schema.Column{
							{
								Name: "string_col",
								Type: arrow.BinaryTypes.String,
							},
						},
					},
				},
			},
			{
				Name:        "relation_table2",
				Description: "Description for second relational table",
				Columns: []schema.Column{
					{
						Name: "string_col",
						Type: arrow.BinaryTypes.String,
					},
				},
			},
		},
	},
	{
		Name:          "incremental_table",
		Description:   "Description for incremental table",
		IsIncremental: true,
		Columns: []schema.Column{
			{
				Name: "int_col",
				Type: arrow.PrimitiveTypes.Int64,
			},
			{
				Name:           "id_col",
				Type:           arrow.PrimitiveTypes.Int64,
				PrimaryKey:     true,
				IncrementalKey: true,
			},
			{
				Name:           "id_col2",
				Type:           arrow.PrimitiveTypes.Int64,
				IncrementalKey: true,
			},
		},
	},
}

func TestGeneratePluginDocs(t *testing.T) {
	p := NewPlugin("test", "v1.0.0", testTables, newTestExecutionClient)

	cup := cupaloy.New(cupaloy.SnapshotSubdirectory("testdata"))

	t.Run("Markdown", func(t *testing.T) {
		tmpdir := t.TempDir()

		err := p.GeneratePluginDocs(tmpdir, "markdown")
		if err != nil {
			t.Fatalf("unexpected error calling GeneratePluginDocs: %v", err)
		}

		expectFiles := []string{"test_table.md", "relation_table.md", "relation_relation_table_a.md", "relation_relation_table_b.md", "incremental_table.md", "README.md"}
		for _, exp := range expectFiles {
			t.Run(exp, func(t *testing.T) {
				output := path.Join(tmpdir, exp)
				got, err := os.ReadFile(output)
				require.NoError(t, err)
				cup.SnapshotT(t, got)
			})
		}
	})

	t.Run("JSON", func(t *testing.T) {
		tmpdir := t.TempDir()

		err := p.GeneratePluginDocs(tmpdir, "json")
		if err != nil {
			t.Fatalf("unexpected error calling GeneratePluginDocs: %v", err)
		}

		expectFiles := []string{"__tables.json"}
		for _, exp := range expectFiles {
			t.Run(exp, func(t *testing.T) {
				output := path.Join(tmpdir, exp)
				got, err := os.ReadFile(output)
				require.NoError(t, err)
				cup.SnapshotT(t, got)
			})
		}
	})
}
