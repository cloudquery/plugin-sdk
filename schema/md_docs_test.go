package schema

import (
	"io/ioutil"
	"os"
	"testing"
)

var testGenMarkdownTables = []*Table{
	{
		Name:        "testGenerateMarkdownTable",
		Description: "testTable description",
		Columns: []Column{
			{
				Name: "testGenerateMarkdownColumn",
				Type: TypeInt,
			},
			{
				Name: "testGenerateMarkdownColumn",
				Type: TypeInt,
			},
		},
		Relations: []*Table{
			{
				Name: "testGenerateMarkdownTableChild",
				Columns: []Column{
					{
						Name: "testGenerateMarkdownColumnChild",
						Type: TypeString,
					},
				},
			},
		},
	},
}

func TestGenerateMarkdownTree(t *testing.T) {
	dir, err := ioutil.TempDir(".", "markdowntest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	if err := GenerateMarkdownTree(testGenMarkdownTables, dir); err != nil {
		t.Fatal(err)
	}
}
