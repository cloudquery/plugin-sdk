package plugins

import (
	_ "embed"
	"os"
	"path"
	"testing"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/google/go-cmp/cmp"
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
		},
	},
}

//go:embed testdata/test_table.md
var testTableMd []byte

//go:embed testdata/relation_table.md
var relationTableMd []byte

//go:embed testdata/README.md
var readmeMd []byte

var expectedDocFiles = map[string][]byte{
	"test_table.md":     testTableMd,
	"relation_table.md": relationTableMd,
	"README.md":         readmeMd,
}

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

	for filename, content := range expectedDocFiles {
		t.Run(filename, func(t *testing.T) {
			output := path.Join(tmpdir, filename)
			got, err := os.ReadFile(output)
			if err != nil {
				t.Fatalf("failed to read file %s: %v", output, err)
			}
			if diff := cmp.Diff(content, got); diff != "" {
				t.Fatalf("unexpected file contents %s: %v", output, diff)
			}
		})
	}
}
