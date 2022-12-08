package schema

import "testing"

var testTable = &Table{
	Name:    "test",
	Columns: []Column{},
	Relations: []*Table{
		{
			Name:    "test2",
			Columns: []Column{},
		},
	},
}

func TestTablesFlatten(t *testing.T) {
	tables := Tables{testTable}.FlattenTables()
	if len(tables) != 2 {
		t.Fatal("expected 2 tables")
	}
}
