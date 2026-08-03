package schema

import (
	"math"
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
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
	tg := NewTestDataGenerator(0)
	_ = tg.Generate(table, GenTestDataOptions{})
}

func countInt64Leaves(arr arrow.Array) (total, inexact int) {
	switch col := arr.(type) {
	case *array.Int64:
		for i := 0; i < col.Len(); i++ {
			if col.IsNull(i) {
				continue
			}
			total++
			if v := col.Value(i); int64(float64(v)) != v {
				inexact++
			}
		}
	case *array.Uint64:
		for i := 0; i < col.Len(); i++ {
			if col.IsNull(i) {
				continue
			}
			total++
			if v := col.Value(i); uint64(float64(v)) != v {
				inexact++
			}
		}
	case *array.Struct:
		for i := 0; i < col.NumField(); i++ {
			leaves, bad := countInt64Leaves(col.Field(i))
			total, inexact = total+leaves, inexact+bad
		}
	case *array.List:
		total, inexact = countInt64Leaves(col.ListValues())
	case *array.LargeList:
		total, inexact = countInt64Leaves(col.ListValues())
	case *array.Map:
		kt, kx := countInt64Leaves(col.Keys())
		it, ix := countInt64Leaves(col.Items())
		total, inexact = kt+it, kx+ix
	}
	return total, inexact
}

const minInt64Leaves = 100

func TestGenTestDataMaxIntegerBits(t *testing.T) {
	table := TestTable("test", TestSourceOptions{})

	for _, tc := range []struct {
		name           string
		maxIntegerBits int
		wantExact      bool
	}{
		{name: "unbounded", maxIntegerBits: 0, wantExact: false},
		{name: "float64 safe", maxIntegerBits: Float64SafeIntegerBits, wantExact: true},
		{name: "one below signed width", maxIntegerBits: 62, wantExact: false},
		{name: "signed width", maxIntegerBits: 63, wantExact: false},
		{name: "full width", maxIntegerBits: 64, wantExact: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tg := NewTestDataGenerator(0)
			record := tg.Generate(table, GenTestDataOptions{MaxRows: 50, MaxIntegerBits: tc.maxIntegerBits})

			total, inexact := 0, 0
			for i := range record.Schema().Fields() {
				leaves, bad := countInt64Leaves(record.Column(i))
				total, inexact = total+leaves, inexact+bad
			}

			if total < minInt64Leaves {
				t.Fatalf("inspected only %d 64-bit integer values, expected at least %d; the walk is not reaching nested columns", total, minInt64Leaves)
			}
			if tc.wantExact && inexact > 0 {
				t.Errorf("MaxIntegerBits=%d generated %d integers that do not survive a float64 round-trip", tc.maxIntegerBits, inexact)
			}
			if !tc.wantExact && inexact == 0 {
				t.Errorf("MaxIntegerBits=%d generated no integers beyond float64 precision", tc.maxIntegerBits)
			}
		})
	}
}

func TestTestTableFixtureHasNestedLargeIntegers(t *testing.T) {
	table := TestTable("test", TestSourceOptions{})
	record := NewTestDataGenerator(0).Generate(table, GenTestDataOptions{MaxRows: 50})

	nested := 0
	for i := range record.Schema().Fields() {
		switch record.Column(i).(type) {
		case *array.Int64, *array.Uint64:
		default:
			_, inexact := countInt64Leaves(record.Column(i))
			nested += inexact
		}
	}

	if nested == 0 {
		t.Fatal("expected nested columns to carry integers beyond float64 precision, so the bounded case proves propagation")
	}
}

func TestGenTestDataMaxIntegerBitsUnsetPreservesFullRange(t *testing.T) {
	if got, want := signedInt64Bound(0), int64(^uint64(0)>>1); got != want {
		t.Errorf("signedInt64Bound(0) = %d, want the full signed range %d", got, want)
	}
	for _, v := range []uint64{0, 1, math.MaxInt64, math.MaxUint64} {
		if got := boundUint64(v, 0); got != v {
			t.Errorf("boundUint64(%d, 0) = %d, want it unchanged", v, got)
		}
	}

	wantInt64 := []int64{-8717895732742165505, -7144924247938981575, -1395437218309923052, -4345851588384648695, -7242748068272024738}
	wantUint64 := []uint64{8717895732742165505, 16368296284793757383, 1395437218309923052, 13569223625239424503, 7242748068272024738}

	table := TestTable("test", TestSourceOptions{})
	record := NewTestDataGenerator(0).Generate(table, GenTestDataOptions{MaxRows: len(wantInt64)})

	for i, field := range record.Schema().Fields() {
		switch field.Name {
		case "int64":
			col := record.Column(i).(*array.Int64)
			for j, want := range wantInt64 {
				if got := col.Value(j); got != want {
					t.Errorf("int64 row %d = %d, want %d", j, got, want)
				}
			}
		case "uint64":
			col := record.Column(i).(*array.Uint64)
			for j, want := range wantUint64 {
				if got := col.Value(j); got != want {
					t.Errorf("uint64 row %d = %d, want %d", j, got, want)
				}
			}
		}
	}
}

func TestSignedInt64BoundIsUsableWithInt63n(t *testing.T) {
	for _, bits := range []int{-1, 0, 1, 52, 53, 62, 63, 64, 65} {
		if got := signedInt64Bound(bits); got <= 0 {
			t.Errorf("signedInt64Bound(%d) = %d, which panics rand.Int63n", bits, got)
		}
	}
}

func TestBoundUint64Width(t *testing.T) {
	for _, tc := range []struct {
		bits int
		want uint64
	}{
		{bits: 1, want: 1},
		{bits: 53, want: 1<<53 - 1},
		{bits: 63, want: 1<<63 - 1},
	} {
		if got := boundUint64(math.MaxUint64, tc.bits); got != tc.want {
			t.Errorf("boundUint64(MaxUint64, %d) = %d, want %d", tc.bits, got, tc.want)
		}
	}
}
