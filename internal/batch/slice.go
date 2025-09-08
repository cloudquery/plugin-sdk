package batch

import (
	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/util"
)

type (
	SlicedRecord struct {
		arrow.RecordBatch
		Bytes       int64 // we need this as the util.TotalRecordSize will report the full size even for the sliced record
		bytesPerRow int64
	}
)

func (s *SlicedRecord) split(limit *Cap) (add *SlicedRecord, toFlush []arrow.RecordBatch, rest *SlicedRecord) {
	if s == nil {
		return nil, nil, nil
	}

	add = s.getAdd(limit)
	if add != nil {
		limit.add(add.Bytes, add.NumRows())
	}

	if s.RecordBatch == nil {
		// all processed
		return add, nil, nil
	}

	toFlush = s.getToFlush(limit)
	if s.RecordBatch == nil {
		// all processed
		return add, toFlush, nil
	}

	// set bytes & rows new values
	limit.set(s.Bytes, s.NumRows())
	return add, toFlush, s
}

func (s *SlicedRecord) getAdd(limit *Cap) *SlicedRecord {
	rowsByBytes := limit.bytes.remainingPerN(s.bytesPerRow)
	rows := limit.rows.remaining()
	switch {
	case rows < 0:
		rows = rowsByBytes
	case rows > rowsByBytes && rowsByBytes >= 0:
		rows = rowsByBytes
	}

	switch {
	case rows == 0:
		return nil
	case rows < 0, rows >= s.NumRows():
		// grab the whole record (either no limits or not overflowing)
		res := *s
		s.Bytes = 0
		s.RecordBatch = nil
		return &res
	}

	res := SlicedRecord{
		RecordBatch: s.NewSlice(0, rows),
		Bytes:       rows * s.bytesPerRow,
		bytesPerRow: s.bytesPerRow,
	}
	s.RecordBatch = s.NewSlice(rows, s.NumRows())
	s.Bytes -= res.Bytes
	return &res
}

func (s *SlicedRecord) getToFlush(limit *Cap) []arrow.RecordBatch {
	rowsByBytes := limit.bytes.capPerN(s.bytesPerRow)
	rows := limit.rows.cap()
	switch {
	case rows < 0:
		rows = rowsByBytes
	case rows > rowsByBytes && rowsByBytes >= 0:
		rows = rowsByBytes
	}

	switch {
	case rows == 0:
		// not even a single row fits
		// we still need to process this, so slice by single row
		return s.slice()
	case rows < 0:
		// as s.Record != nil we know that the limits are there in place & the s.Record.NumRows() > 0
		panic("should never be here")
	case rows > s.NumRows():
		// no need to flush anything, as the amount of rows isn't enough to grant this
		return nil
	}

	flush := make([]arrow.RecordBatch, 0, s.NumRows()/rows)
	offset := int64(0)
	for offset+rows <= s.NumRows() {
		flush = append(flush, s.NewSlice(offset, offset+rows))
		offset += rows
	}
	if offset == s.NumRows() {
		// we processed everything for flush
		s.RecordBatch = nil
		s.Bytes = 0
		return flush
	}

	// set record to the remainder
	s.RecordBatch = s.NewSlice(offset, s.NumRows())
	s.Bytes = s.NumRows() * s.bytesPerRow

	return flush
}

func (s *SlicedRecord) slice() []arrow.RecordBatch {
	res := make([]arrow.RecordBatch, s.NumRows())
	for i := int64(0); i < s.NumRows(); i++ {
		res[i] = s.NewSlice(i, i+1)
	}
	return res
}

func newSlicedRecord(r arrow.RecordBatch) *SlicedRecord {
	if r.NumRows() == 0 {
		return nil
	}
	res := SlicedRecord{
		RecordBatch: r,
		Bytes:       util.TotalRecordSize(r),
	}
	res.bytesPerRow = res.Bytes / r.NumRows()
	return &res
}

// SliceRecord will return the SlicedRecord you can add to the batch given the restrictions provided (if any).
// The meaning of the returned values:
// - `add` is good to be added to the current batch that the caller is assembling
// - `flush` represents sliced arrow.RecordBatch that needs own batch to be flushed
// - `remaining` represents the overflow of the batch after `add` & `flush` are processed
// Note that the `limit` provided will not be updated.
func SliceRecord(r arrow.RecordBatch, limit *Cap) (add *SlicedRecord, flush []arrow.RecordBatch, remaining *SlicedRecord) {
	l := *limit // copy value
	return newSlicedRecord(r).split(&l)
}
