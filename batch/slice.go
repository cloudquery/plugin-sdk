package batch

import (
	"github.com/apache/arrow/go/v16/arrow"
	"github.com/apache/arrow/go/v16/arrow/util"
)

type SlicedRecord struct {
	arrow.Record
	Bytes int64 // we need this as the util.TotalRecordSize will report the full size even for the sliced record
}

// SliceRecord will return the SlicedRecord you can add to the batch given the restrictions provided (if any).
// The meaning of the returned values:
// - `add` is good to be added to the current batch that the caller is assembling
// - `flush` represents sliced arrow.Record that needs own batch to be flushed
// - `remaining` represents the overflow of the batch after `add` & `flush` are processed
func SliceRecord(r arrow.Record, bytes, rows *Cap[int64]) (add *SlicedRecord, flush []SlicedRecord, remaining *SlicedRecord) {
	util.TotalRecordSize()
}
