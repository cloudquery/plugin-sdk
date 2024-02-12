package schema

import "testing"

func TestTestSourceColumns_Default(t *testing.T) {
	// basic sanity check for tested columns
	table := TestTable("test", TestSourceOptions{})
	if len(table.Columns) < 73 {
		t.Fatalf("expected at least 73 columns by default got: %d ", len(table.Columns))
	}
	// test some specific columns
	checkColumnsExist(t, table.Columns, []string{"int64", "date32", "timestamp_us", "string", "struct", "string_list"})
}

func TestTestSourceColumns_SkipAll(t *testing.T) {
	table := TestTable("test", TestSourceOptions{
		SkipLists:      true,
		SkipTimestamps: true,
		SkipDates:      true,
		SkipMaps:       true,
		SkipStructs:    true,
		SkipIntervals:  true,
		SkipDurations:  true,
		SkipTimes:      true,
		SkipLargeTypes: true,
	})

	// test some specific columns
	checkColumnsExist(t, table.Columns, []string{"int64", "string"})
	checkColumnsDontExist(t, table.Columns, []string{"date32", "struct", "string_map"})
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

func TestGenTestData(*testing.T) {
	table := TestTable("test", TestSourceOptions{})
	// smoke test that no panics
	tg := NewTestDataGenerator(0)
	_ = tg.Generate(table, GenTestDataOptions{})
}
