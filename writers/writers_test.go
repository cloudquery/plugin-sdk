package writers_test

import (
	"context"
	"math/rand"
	"runtime"
	"sort"
	"strconv"
	"testing"

	"github.com/apache/arrow/go/v16/arrow"
	"github.com/apache/arrow/go/v16/arrow/array"
	"github.com/apache/arrow/go/v16/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/writers"
	"github.com/cloudquery/plugin-sdk/v4/writers/batchwriter"
	"github.com/cloudquery/plugin-sdk/v4/writers/mixedbatchwriter"
	"github.com/cloudquery/plugin-sdk/v4/writers/streamingbatchwriter"
	"golang.org/x/exp/maps"
)

type bCase struct {
	name string
	wr   writers.Writer
	rec  func() arrow.Record
}

func BenchmarkWriterMemory(b *testing.B) {
	batchwriterOpts := map[string][]batchwriter.Option{
		"defaults":           nil,
		"batch10k bytes100M": {batchwriter.WithBatchSizeBytes(100000000), batchwriter.WithBatchSize(10000)},
	}
	mixedbatchwriterOpts := map[string][]mixedbatchwriter.Option{
		"defaults":           nil,
		"batch10k bytes100M": {mixedbatchwriter.WithBatchSizeBytes(100000000), mixedbatchwriter.WithBatchSize(10000)},
	}
	streamingbatchwriterOpts := map[string][]streamingbatchwriter.Option{
		"defaults":  nil,
		"bytes100M": {streamingbatchwriter.WithBatchSizeBytes(100000000)},
	}

	var bCases []bCase
	bCases = append(bCases, writerMatrix("BatchWriter", batchwriter.New, newBatchWriterClient(), makeRecord, batchwriterOpts)...)
	bCases = append(bCases, writerMatrix("BatchWriter wide", batchwriter.New, newBatchWriterClient(), makeWideRecord, batchwriterOpts)...)
	bCases = append(bCases, writerMatrix("MixedBatchWriter", mixedbatchwriter.New, newMixedBatchWriterClient(), makeRecord, mixedbatchwriterOpts)...)
	bCases = append(bCases, writerMatrix("MixedBatchWriter wide", mixedbatchwriter.New, newMixedBatchWriterClient(), makeWideRecord, mixedbatchwriterOpts)...)
	bCases = append(bCases, writerMatrix("StreamingBatchWriter", streamingbatchwriter.New, newStreamingBatchWriterClient(), makeRecord, streamingbatchwriterOpts)...)
	bCases = append(bCases, writerMatrix("StreamingBatchWriter wide", streamingbatchwriter.New, newStreamingBatchWriterClient(), makeWideRecord, streamingbatchwriterOpts)...)

	for _, c := range bCases {
		c := c
		b.Run(c.name, func(b *testing.B) {
			var (
				mStart runtime.MemStats
				mEnd   runtime.MemStats
			)

			ch := make(chan message.WriteMessage)
			errCh := make(chan error)
			go func() {
				defer close(errCh)
				errCh <- c.wr.Write(context.Background(), ch)
			}()

			runtime.ReadMemStats(&mStart)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				rec := c.rec()
				ch <- &message.WriteInsert{
					Record: rec,
				}
			}
			close(ch)
			err := <-errCh

			b.StopTimer()

			if err != nil {
				b.Fatal(err)
			}

			runtime.ReadMemStats(&mEnd)

			allocatedBytes := mEnd.Alloc - mStart.Alloc
			b.ReportMetric(float64(allocatedBytes)/float64(b.N), "bytes/op") // this is different from -benchmem result "B/op"
		})
	}
}

func makeRecord() func() arrow.Record {
	table := &schema.Table{
		Name: "test_table",
		Columns: schema.ColumnList{
			{
				Name: "col1",
				Type: arrow.BinaryTypes.String,
			},
		},
	}
	sc := table.ToArrowSchema()

	return func() arrow.Record {
		bldr := array.NewRecordBuilder(memory.DefaultAllocator, sc)
		bldr.Field(0).(*array.StringBuilder).Append("test")
		return bldr.NewRecord()
	}
}

func makeWideRecord() func() arrow.Record {
	table := &schema.Table{
		Name: "test_wide_table",
		Columns: schema.ColumnList{
			{
				Name: "col1",
				Type: arrow.BinaryTypes.String,
			},
		},
	}

	const numWideCols = 200
	randVals := make([]int64, numWideCols)
	for i := 0; i < numWideCols; i++ {
		table.Columns = append(table.Columns, schema.Column{
			Name: "wide_col" + strconv.Itoa(i),
			Type: arrow.PrimitiveTypes.Int64,
		})
		randVals[i] = rand.Int63()
	}
	sc := table.ToArrowSchema()

	return func() arrow.Record {
		bldr := array.NewRecordBuilder(memory.DefaultAllocator, sc)
		bldr.Field(0).(*array.StringBuilder).Append("test")
		for i := 0; i < numWideCols; i++ {
			bldr.Field(i + 1).(*array.Int64Builder).Append(randVals[i])
		}
		return bldr.NewRecord()
	}
}

func writerMatrix[T writers.Writer, C any, O ~func(T)](prefix string, constructor func(C, ...O) (T, error), client C, recordMaker func() func() arrow.Record, optsMatrix map[string][]O) []bCase {
	bCases := make([]bCase, 0, len(optsMatrix))

	k := maps.Keys(optsMatrix)
	sort.Strings(k)

	for _, name := range k {
		opts := optsMatrix[name]
		wr, err := constructor(client, opts...)
		if err != nil {
			panic(err)
		}
		bCases = append(bCases, bCase{
			name: prefix + " " + name,
			wr:   wr,
			rec:  recordMaker(),
		})
	}
	return bCases
}

type mixedbatchwriterClient struct {
	mixedbatchwriter.IgnoreMigrateTableBatch
	mixedbatchwriter.UnimplementedDeleteStaleBatch
	mixedbatchwriter.UnimplementedDeleteRecordsBatch
}

func newMixedBatchWriterClient() mixedbatchwriter.Client {
	return &mixedbatchwriterClient{}
}

func (mixedbatchwriterClient) InsertBatch(_ context.Context, msgs message.WriteInserts) error {
	for _, m := range msgs {
		m.Record.Release()
	}
	return nil
}

var _ mixedbatchwriter.Client = (*mixedbatchwriterClient)(nil)

type batchwriterClient struct {
	batchwriter.IgnoreMigrateTables
	batchwriter.UnimplementedDeleteStale
	batchwriter.UnimplementedDeleteRecord
}

func newBatchWriterClient() batchwriter.Client {
	return &batchwriterClient{}
}

func (batchwriterClient) WriteTableBatch(_ context.Context, _ string, msgs message.WriteInserts) error {
	for _, m := range msgs {
		m.Record.Release()
	}
	return nil
}

var _ batchwriter.Client = (*batchwriterClient)(nil)

type streamingbatchwriterClient struct {
	streamingbatchwriter.IgnoreMigrateTable
	streamingbatchwriter.UnimplementedDeleteStale
	streamingbatchwriter.UnimplementedDeleteRecords
}

func newStreamingBatchWriterClient() streamingbatchwriter.Client {
	return &streamingbatchwriterClient{}
}

func (streamingbatchwriterClient) WriteTable(_ context.Context, ch <-chan *message.WriteInsert) error {
	for m := range ch {
		m.Record.Release()
	}
	return nil
}

var _ streamingbatchwriter.Client = (*streamingbatchwriterClient)(nil)
