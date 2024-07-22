package batch

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/apache/arrow/go/v17/arrow/util"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/stretchr/testify/assert"
)

func TestSliceRecord(t *testing.T) {
	for run := 0; run < 5; run++ {
		rows := rand.Intn(100) + 5
		t.Run(strconv.Itoa(rows), func(t *testing.T) {
			t.Parallel()
			table := schema.TestTable(fmt.Sprintf("test_%d_rows", rows), schema.TestSourceOptions{})
			tg := schema.NewTestDataGenerator(0)
			record := tg.Generate(table, schema.GenTestDataOptions{
				MaxRows:    rows,
				SourceName: "test",
				SyncTime:   time.Now(),
			})

			recordRows, recordBytes := record.NumRows(), util.TotalRecordSize(record)

			t.Run("only add", func(t *testing.T) {
				limit := CappedAt(0, 0)
				add, toFlush, rest := SliceRecord(record, limit)
				assert.Equal(t, limit.bytes, capped{})
				assert.Equal(t, limit.rows, capped{})
				assert.NotNil(t, add)
				assert.Equal(t, recordRows, add.NumRows())
				assert.Equal(t, recordBytes, add.Bytes)
				assert.Empty(t, toFlush)
				assert.Nil(t, rest)
			})

			t.Run("only single toFlush", func(t *testing.T) {
				limit := &Cap{
					bytes: capped{current: recordBytes, limit: recordBytes},
					rows:  capped{current: recordRows, limit: recordRows},
				}
				add, toFlush, rest := SliceRecord(record, limit)
				assert.Equal(t, limit.bytes, capped{current: recordBytes, limit: recordBytes})
				assert.Equal(t, limit.rows, capped{current: recordRows, limit: recordRows})
				assert.Nil(t, add)
				assert.NotEmpty(t, toFlush)
				assert.Len(t, toFlush, 1)
				r := toFlush[0]
				assert.Equal(t, recordRows, r.NumRows())
				assert.Nil(t, rest)
			})

			t.Run("full - by rows", func(t *testing.T) {
				limit := &Cap{rows: capped{current: recordRows / 10, limit: recordRows / 5}}
				remaining := recordRows

				add, toFlush, rest := SliceRecord(record, limit)
				assert.Equal(t, limit.bytes, capped{})
				assert.Equal(t, limit.rows, capped{current: recordRows / 10, limit: recordRows / 5})

				// if we could add some rows
				if (recordRows/5)-(recordRows/10) > 0 {
					assert.NotNil(t, add)
					assert.LessOrEqual(t, add.NumRows(), recordRows/5)
					assert.LessOrEqual(t, add.Bytes, recordBytes/5)
					remaining -= add.NumRows()
				} else {
					assert.Nil(t, add)
				}

				assert.NotEmpty(t, toFlush)
				assert.GreaterOrEqual(t, len(toFlush), 4)
				for _, f := range toFlush {
					assert.LessOrEqual(t, f.NumRows(), recordRows/5)
					remaining -= f.NumRows()
				}

				assert.GreaterOrEqual(t, remaining, int64(0))
				if remaining == 0 {
					assert.Nil(t, rest)
					return
				}

				assert.NotNil(t, rest)
				assert.Less(t, remaining, recordRows/5)
				assert.Equal(t, remaining, rest.NumRows())
				assert.Less(t, rest.Bytes, recordBytes/5)
			})

			t.Run("full - by bytes", func(t *testing.T) {
				limit := &Cap{bytes: capped{current: recordBytes / 10, limit: recordBytes / 5}}
				remaining := recordRows

				add, toFlush, rest := SliceRecord(record, limit)
				assert.Equal(t, limit.bytes, capped{current: recordBytes / 10, limit: recordBytes / 5})
				assert.Equal(t, limit.rows, capped{})

				// if we could add some rows
				if (recordBytes/5)-(recordBytes/10) >= util.TotalRecordSize(record)/record.NumRows() {
					assert.NotNil(t, add)
					assert.LessOrEqual(t, add.NumRows(), recordRows/5)
					assert.LessOrEqual(t, add.Bytes, recordBytes/5)
					remaining -= add.NumRows()
				} else {
					assert.Nil(t, add)
				}

				assert.NotEmpty(t, toFlush)
				assert.GreaterOrEqual(t, len(toFlush), 4)
				for _, f := range toFlush {
					assert.LessOrEqual(t, f.NumRows(), recordRows/5)
					remaining -= f.NumRows()
				}

				assert.GreaterOrEqual(t, remaining, int64(0))
				if remaining == 0 {
					assert.Nil(t, rest)
					return
				}

				assert.NotNil(t, rest)
				assert.Less(t, remaining, recordRows/5)
				assert.Equal(t, remaining, rest.NumRows())
				assert.Less(t, rest.Bytes, recordBytes/5)
			})
		})
	}
}
