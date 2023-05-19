package schema

import "testing"

func TestTestSourceColumns_Default(t *testing.T) {
	// basic sanity check for tested columns
	defaults := TestSourceColumns()
	if len(defaults) < 73 {
		t.Fatalf("expected at least 73 columns by default got: %d ", len(defaults))
	}
	// test some specific columns
	checkColumnsExist(t, defaults, []string{"int64", "date32", "timestamp_us", "string", "struct", "string_list"})
}

func TestTestSourceColumns_SkipAll(t *testing.T) {
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
	// test some specific columns
	checkColumnsExist(t, skipAll, []string{"int64", "timestamp_us", "string", "string_list"})
	checkColumnsDontExist(t, skipAll, []string{"date32", "struct", "string_map"})
}

func checkColumnsExist(t *testing.T, list ColumnList, cols []string) {
	for _, col := range cols {
		if list.Get(col) == nil {
			t.Errorf("expected column %s to be present", col)
		}
	}
}

func checkColumnsDontExist(t *testing.T, list ColumnList, cols []string) {
	for _, col := range cols {
		if list.Get(col) != nil {
			t.Errorf("expected no %s column", col)
		}
	}
}
