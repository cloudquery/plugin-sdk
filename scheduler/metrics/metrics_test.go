package metrics

import "testing"

func TestMetrics(t *testing.T) {
	s := NewMetrics("test_invocation_id")
	s.TableClient["test_table"] = make(map[string]*tableClientMetrics)
	s.TableClient["test_table"]["testExecutionClient"] = &tableClientMetrics{
		resources: 1,
		errors:    2,
		panics:    3,
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

	other := NewMetrics("test_invocation_id")
	other.TableClient["test_table"] = make(map[string]*tableClientMetrics)
	other.TableClient["test_table"]["testExecutionClient"] = &tableClientMetrics{
		resources: 1,
		errors:    2,
		panics:    3,
	}
	if !s.Equal(other) {
		t.Fatal("expected metrics to be equal")
	}
}
