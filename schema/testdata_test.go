package schema

import "testing"

func TestTestSourceColumns(t *testing.T) {
	// basic sanity check for tested columns
	defaults := TestSourceColumns()
	if len(defaults) < 100 {
		t.Fatal("expected at least 100 columns by default")
	}
	skipAll := ColumnList(TestSourceColumns(
		WithTestSourceSkipStructs(),
		WithTestSourceSkipMaps(),
		WithTestSourceSkipDates(),
		WithTestSourceSkipTimes(),
		WithTestSourceSkipTimestamps(),
		WithTestSourceSkipDurations(),
		WithTestSourceSkipIntervals(),
		WithTestSourceSkipLargeTypes(),
	))
	if skipAll.Get("struct") != nil {
		t.Fatal("expected no structs when WithTestSourceSkipStructs is used")
	}
	if skipAll.Get("map") != nil {
		t.Fatal("expected no maps when WithTestSourceSkipMaps is used")
	}
	if skipAll.Get("date32") != nil {
		t.Fatal("expected no date32 when WithTestSourceSkipDates is used")
	}
	if skipAll.Get("time32") != nil {
		t.Fatal("expected no times when WithTestSourceSkipTimes is used")
	}
	if skipAll.Get("timestamp_us") == nil {
		t.Fatal("expected no microsecond timestamps even when WithTestSourceSkipTimestamps is used")
	}
	if skipAll.Get("timestamp_ns") != nil {
		t.Fatal("expected no nansecond timestamps when WithTestSourceSkipTimestamps is used")
	}
	if skipAll.Get("monthinterval") != nil {
		t.Fatal("expected no interval when WithTestSourceSkipIntervals is used")
	}
}
