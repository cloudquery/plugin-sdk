package batch

import (
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/apache/arrow-go/v18/arrow/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func boolRecord(t *testing.T, rows int) arrow.RecordBatch {
	t.Helper()

	sc := arrow.NewSchema([]arrow.Field{{Name: "b", Type: arrow.FixedWidthTypes.Boolean}}, nil)
	builder := array.NewRecordBuilder(memory.DefaultAllocator, sc)
	defer builder.Release()

	for range rows {
		builder.Field(0).(*array.BooleanBuilder).Append(true)
	}
	return builder.NewRecordBatch()
}

func TestSlicedRecordBytesNeverNegative(t *testing.T) {
	const rows = 20000
	record := boolRecord(t, rows)
	t.Logf("real size: %d bytes for %d rows", util.TotalRecordSize(record), record.NumRows())

	limit := CappedAt(18000, 0)
	limit.bytes.current = 4000

	add, toFlush, rest := SliceRecord(record, limit)

	require.NotNil(t, add)
	assert.GreaterOrEqual(t, add.Bytes, int64(0), "add.Bytes must not be negative")
	if rest != nil {
		t.Logf("rest: %d rows, %d bytes", rest.NumRows(), rest.Bytes)
		assert.GreaterOrEqual(t, rest.Bytes, int64(0), "rest.Bytes must not be negative")
	}

	var total int64
	if add != nil {
		total += add.NumRows()
	}
	for _, r := range toFlush {
		total += r.NumRows()
	}
	if rest != nil {
		total += rest.NumRows()
	}
	assert.Equal(t, int64(rows), total, "no rows may be lost")
}
