package plugin

import (
	"strings"
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
)

func int64Record(t *testing.T, values []int64) arrow.RecordBatch {
	t.Helper()
	sc := arrow.NewSchema([]arrow.Field{{Name: "int64", Type: arrow.PrimitiveTypes.Int64}}, nil)
	bldr := array.NewRecordBuilder(memory.DefaultAllocator, sc)
	defer bldr.Release()
	bldr.Field(0).(*array.Int64Builder).AppendValues(values, nil)
	return bldr.NewRecordBatch()
}

func TestRecordsDiffHintsAtMaxIntegerBits(t *testing.T) {
	sc := arrow.NewSchema([]arrow.Field{{Name: "int64", Type: arrow.PrimitiveTypes.Int64}}, nil)
	exact := int64Record(t, []int64{-8717895732742165505})
	rounded := int64Record(t, []int64{-8717895732742165504})

	diff := RecordsDiff(sc, []arrow.RecordBatch{rounded}, []arrow.RecordBatch{exact})
	if diff == "" {
		t.Fatal("expected a diff between an exact and a float64-rounded int64")
	}
	if !strings.Contains(diff, "MaxIntegerBits") {
		t.Errorf("expected the diff to point at MaxIntegerBits, got: %s", diff)
	}
}

func TestRecordsDiffOmitsHintForUnrelatedDifference(t *testing.T) {
	sc := arrow.NewSchema([]arrow.Field{{Name: "int64", Type: arrow.PrimitiveTypes.Int64}}, nil)
	want := int64Record(t, []int64{1})
	have := int64Record(t, []int64{2})

	diff := RecordsDiff(sc, []arrow.RecordBatch{have}, []arrow.RecordBatch{want})
	if diff == "" {
		t.Fatal("expected a diff between 1 and 2")
	}
	if strings.Contains(diff, "MaxIntegerBits") {
		t.Errorf("did not expect a MaxIntegerBits hint for a difference float64 cannot explain, got: %s", diff)
	}
}
