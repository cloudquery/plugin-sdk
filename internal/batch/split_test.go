package batch

import (
	"testing"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/apache/arrow-go/v18/arrow/util"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRecord(t *testing.T, rows int) arrow.RecordBatch {
	t.Helper()

	table := schema.TestTable("test_table", schema.TestSourceOptions{})
	return schema.NewTestDataGenerator(0).Generate(table, schema.GenTestDataOptions{
		MaxRows:    rows,
		SourceName: "test",
		SyncTime:   time.Now(),
	})
}

func totalRows(records []arrow.RecordBatch) int64 {
	var rows int64
	for _, record := range records {
		rows += record.NumRows()
	}
	return rows
}

func TestSplitRecord(t *testing.T) {
	const rows = 20
	record := testRecord(t, rows)
	bytesPerRow := util.TotalRecordSize(record) / record.NumRows()

	t.Run("no limit", func(t *testing.T) {
		records := SplitRecord(record, CappedAt(0, 0))
		require.Len(t, records, 1)
		assert.Equal(t, int64(rows), records[0].NumRows())
	})

	t.Run("limit above record size", func(t *testing.T) {
		records := SplitRecord(record, CappedAt(0, 2*rows))
		require.Len(t, records, 1)
		assert.Equal(t, int64(rows), records[0].NumRows())
	})

	t.Run("rows limit", func(t *testing.T) {
		records := SplitRecord(record, CappedAt(0, 6))
		assert.Equal(t, int64(rows), totalRows(records))
		for _, r := range records {
			assert.LessOrEqual(t, r.NumRows(), int64(6))
		}
	})

	t.Run("bytes limit", func(t *testing.T) {
		records := SplitRecord(record, CappedAt(5*bytesPerRow, 0))
		assert.Equal(t, int64(rows), totalRows(records))
		for _, r := range records {
			assert.LessOrEqual(t, r.NumRows(), int64(5))
		}
	})

	t.Run("no single row fits", func(t *testing.T) {
		records := SplitRecord(record, CappedAt(1, 0))
		assert.Len(t, records, int(record.NumRows()), "each row gets a record of its own")
		assert.Equal(t, record.NumRows(), totalRows(records))
	})

	t.Run("empty record", func(t *testing.T) {
		empty := testRecord(t, 0)
		records := SplitRecord(empty, CappedAt(1, 1))
		require.Len(t, records, 1)
		assert.Equal(t, int64(0), records[0].NumRows())
	})
}

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
