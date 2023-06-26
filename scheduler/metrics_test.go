package scheduler

import "testing"

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
