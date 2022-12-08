package schema

import "testing"

var testTable = &Table{
	Name:    "test",
	Columns: []Column{},
	Relations: []*Table{
		{
			Name:    "test_sub",
			Columns: []Column{},
		},
	},
}

var testTable2 = &Table{
	Name:    "test2",
	Columns: []Column{},
	Relations: []*Table{
		{
			Name:    "test2_sub",
			Columns: []Column{},
			Relations: []*Table{
				{
					Name:    "test2_sub_sub",
					Columns: []Column{},
				},
			},
		},
	},
}

func TestTablesFlatten(t *testing.T) {
	tables := Tables{testTable}.FlattenTables()
	if len(tables) != 2 {
		t.Fatal("expected 2 tables")
	}
}
