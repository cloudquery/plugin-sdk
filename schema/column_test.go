package schema

import (
	"encoding/json"
	"testing"
)

func TestColumnListMarshal(t *testing.T) {
	c := ColumnList{
		{
			Name: "test",
			Type: TypeBool,
		},
	}
	b, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	var res ColumnList
	if err := json.Unmarshal(b, &res); err != nil {
		t.Fatal(err)
	}
	if len(c) != len(res) {
		t.Fatalf("expected %d columns but got %d", len(c), len(res))
	}
}

func TestColumnListMarshalUnknown(t *testing.T) {
	c := ColumnList{
		{
			Name: "test",
			Type: TypeBool,
		},
		{
			Name: "unknown",
			Type: 10010,
		},
	}
	b, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	var res ColumnList
	if err := json.Unmarshal(b, &res); err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 {
		t.Fatalf("expected 1 columns but got %d", len(res))
	}
}
