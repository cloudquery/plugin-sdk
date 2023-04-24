package testdata

import (
	"testing"

	"github.com/cloudquery/plugin-sdk/v2/schema"
)

func TestTestSourceSchema(t *testing.T) {
	s := TestSourceSchema("test", TestSourceOptions{})
	if schema.TableName(s) != "test" {
		t.Fatal("wrong name")
	}

	for _, f := range s.Fields() {
		t.Log("field:", f)
	}

	_, found := s.FieldsByName(schema.CqSourceNameColumn.Name)
	if !found {
		t.Errorf("_cq_source_name column not found")
	}

	_, found = s.FieldsByName(schema.CqSyncTimeColumn.Name)
	if !found {
		t.Errorf("_cq_sync_time column not found")
	}
}
