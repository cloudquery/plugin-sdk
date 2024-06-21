package batch

import (
	"github.com/apache/arrow/go/v16/arrow"
	"github.com/apache/arrow/go/v16/arrow/util"
)

type (
	SlicedRecord struct {
		arrow.Record
		Bytes       int64 // we need this as the util.TotalRecordSize will report the full size even for the sliced record
		bytesPerRow int64
	}
)

func (s *SlicedRecord) split(limit *Cap) (add *SlicedRecord, toFlush []arrow.Record, rest *SlicedRecord) {
	if s == nil {
		return nil, nil, nil
	}

	add = s.getAdd(limit)
	if add != nil {
		limit.add(add.Bytes, add.NumRows())
	}

	if s.Record == nil {
		// all processed
		return add, nil, nil
	}

	toFlush = s.getToFlush(limit)
	if s.Record == nil {
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
		s.Record = nil
		return &res
	}

	res := SlicedRecord{
		Record:      s.NewSlice(0, rows),
		Bytes:       rows * s.bytesPerRow,
		bytesPerRow: s.bytesPerRow,
	}
	s.Record = s.NewSlice(rows, s.NumRows())
	s.Bytes -= res.Bytes
	return &res
}

func (s *SlicedRecord) getToFlush(limit *Cap) []arrow.Record {
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

	flush := make([]arrow.Record, 0, s.NumRows()/rows)
	offset := int64(0)
	for offset+rows <= s.NumRows() {
		flush = append(flush, s.NewSlice(offset, offset+rows))
		offset += rows
	}
	if offset == s.NumRows() {
		// we processed everything for flush
		s.Record = nil
		s.Bytes = 0
		return flush
	}

	// set record to the remainder
	s.Record = s.NewSlice(offset, s.NumRows())
	s.Bytes = s.NumRows() * s.bytesPerRow

	return flush
}

func (s *SlicedRecord) slice() []arrow.Record {
	res := make([]arrow.Record, s.NumRows())
	for i := int64(0); i < s.NumRows(); i++ {
		res[i] = s.NewSlice(i, i+1)
	}
	return res
}

func newSlicedRecord(r arrow.Record) *SlicedRecord {
	if r.NumRows() == 0 {
		return nil
	}
	res := SlicedRecord{
		Record: r,
		Bytes:  util.TotalRecordSize(r),
	}
	res.bytesPerRow = res.Bytes / r.NumRows()
	return &res
}

// SliceRecord will return the SlicedRecord you can add to the batch given the restrictions provided (if any).
// The meaning of the returned values:
// - `add` is good to be added to the current batch that the caller is assembling
// - `flush` represents sliced arrow.Record that needs own batch to be flushed
// - `remaining` represents the overflow of the batch after `add` & `flush` are processed
// Note that the `limit` provided will not be updated.
func SliceRecord(r arrow.Record, limit *Cap) (add *SlicedRecord, flush []arrow.Record, remaining *SlicedRecord) {
	l := *limit // copy value
	return newSlicedRecord(r).split(&l)
}
