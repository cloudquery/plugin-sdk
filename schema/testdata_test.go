package schema

import (
	"testing"

	"github.com/apache/arrow/go/v13/arrow"
	pb "github.com/cloudquery/plugin-pb-go/pb/plugin/v3"
)

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
	tg := NewTestDataGenerator()
	_ = tg.Generate(table, GenTestDataOptions{})
}

func BenchmarkMultiRow(b *testing.B) {
	table := TestTable("test", TestSourceOptions{})
	tg := NewTestDataGenerator()
	record := tg.Generate(table, GenTestDataOptions{
		SourceName: "test",
		MaxRows:    b.N,
	})
	b.ResetTimer()
	_, err := pb.RecordToBytes(record)
	if err != nil {
		b.Fatal(err)
	}
	b.StopTimer()
}

func BenchmarkSingleRow(b *testing.B) {
	table := TestTable("test", TestSourceOptions{})
	tg := NewTestDataGenerator()
	records := split(tg.Generate(table, GenTestDataOptions{
		SourceName: "test",
		MaxRows:    b.N,
	}))
	b.ResetTimer()
	for i := range records {
		_, err := pb.RecordToBytes(records[i])
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
}

func split(r arrow.Record) []arrow.Record {
	res := make([]arrow.Record, r.NumRows())
	for i := int64(0); i < r.NumRows(); i++ {
		res[i] = r.NewSlice(i, i+1)
	}
	return res
}
