package source

import (
	"testing"
	"time"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/stretchr/testify/assert"
)

func TestMetrics(t *testing.T) {
	s := &Metrics{
		TableClient: make(map[string]map[string]*TableClientMetrics),
	}
	s.TableClient["test_table"] = make(map[string]*TableClientMetrics)
	s.TableClient["test_table"]["testExecutionClient"] = &TableClientMetrics{
		Resources: 1,
		Errors:    2,
		Panics:    3,
	}
	if s.TotalResources() != 1 {
		t.Fatal("expected 1 resource")
	}
	if s.TotalErrors() != 2 {
		t.Fatal("expected 2 error")
	}
	if s.TotalPanics() != 3 {
		t.Fatal("expected 3 panics")
	}

	other := &Metrics{
		TableClient: make(map[string]map[string]*TableClientMetrics),
	}
	other.TableClient["test_table"] = make(map[string]*TableClientMetrics)
	other.TableClient["test_table"]["testExecutionClient"] = &TableClientMetrics{
		Resources: 1,
		Errors:    2,
		Panics:    3,
	}
	if !s.Equal(other) {
		t.Fatal("expected metrics to be equal")
	}
}

func TestInProgressTables(t *testing.T) {
	s := &Metrics{
		TableClient: make(map[string]map[string]*TableClientMetrics),
	}
	s.TableClient["test_table_done"] = make(map[string]*TableClientMetrics)
	s.TableClient["test_table_done"]["testExecutionClient"] = &TableClientMetrics{
		Resources: 1,
		Errors:    2,
		Panics:    3,
		startTime: time.Now(),
		endTime:   time.Now().Add(time.Second),
	}

	s.TableClient["test_table_running1"] = make(map[string]*TableClientMetrics)
	s.TableClient["test_table_running1"]["testExecutionClient"] = &TableClientMetrics{
		Resources: 1,
		Errors:    2,
		Panics:    3,
		startTime: time.Now(),
	}

	s.TableClient["test_table_running2"] = make(map[string]*TableClientMetrics)
	s.TableClient["test_table_running2"]["testExecutionClient"] = &TableClientMetrics{
		Resources: 1,
		Errors:    2,
		Panics:    3,
		startTime: time.Now(),
	}

	assert.ElementsMatch(t, []string{"test_table_running1", "test_table_running2"}, s.InProgressTables())
}

type MockClientMeta struct {
}

func (*MockClientMeta) ID() string {
	return "id"
}

var exampleTableSchema = &schema.Table{
	Name: "toplevel",
	Columns: schema.ColumnList{
		{
			Name: "col1",
			Type: schema.TypeInt,
		},
	},
	Relations: []*schema.Table{
		{
			Name: "child",
			Columns: schema.ColumnList{
				{
					Name: "col1",
					Type: schema.TypeInt,
				},
			},
		},
	},
}

// When a top-level table is marked as done, all child tables should be marked as done as well.
// For child-tables, only the specified table should be marked as done.
func TestMarkEndChildTableNotRecursive(t *testing.T) {
	mockClientMeta := &MockClientMeta{}

	metrics := &Metrics{
		TableClient: make(map[string]map[string]*TableClientMetrics),
	}
	metrics.TableClient["toplevel"] = nil
	metrics.TableClient["child"] = nil

	parentTable := exampleTableSchema
	childTable := exampleTableSchema.Relations[0]

	metrics.initWithClients(parentTable, []schema.ClientMeta{mockClientMeta})
	metrics.MarkStart(parentTable, mockClientMeta.ID())
	metrics.MarkStart(childTable, mockClientMeta.ID())

	assert.ElementsMatch(t, []string{"toplevel", "child"}, metrics.InProgressTables())

	metrics.MarkEnd(childTable, mockClientMeta.ID())

	assert.ElementsMatch(t, []string{"toplevel"}, metrics.InProgressTables())
}

func TestMarkEndTopLevelTableRecursive(t *testing.T) {
	mockClientMeta := &MockClientMeta{}

	metrics := &Metrics{
		TableClient: make(map[string]map[string]*TableClientMetrics),
	}
	metrics.TableClient["toplevel"] = nil
	metrics.TableClient["child"] = nil

	parentTable := exampleTableSchema
	childTable := exampleTableSchema.Relations[0]

	metrics.initWithClients(parentTable, []schema.ClientMeta{mockClientMeta})
	metrics.MarkStart(parentTable, mockClientMeta.ID())
	metrics.MarkStart(childTable, mockClientMeta.ID())

	assert.ElementsMatch(t, []string{"toplevel", "child"}, metrics.InProgressTables())

	metrics.MarkEnd(parentTable, mockClientMeta.ID())

	assert.Empty(t, metrics.InProgressTables())
}
