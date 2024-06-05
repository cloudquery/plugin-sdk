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
	grabBytesRows := int64(-1)
	if limit.bytes.limit > 0 {
		grabBytesRows = (limit.bytes.limit - limit.bytes.current) / s.bytesPerRow
	}

	grabRows := int64(-1)
	if limit.rows.limit > 0 {
		grabRows = limit.rows.limit - limit.rows.current
	}

	if grabRows < 0 && grabBytesRows < 0 {
		// no limits
		res := *s
		s.Bytes = 0
		s.Record = nil
		return &res
	}

	grabRows = min(max(grabRows, 0), max(grabBytesRows, 0))
	if grabRows == 0 {
		return nil
	}
	if grabRows >= s.NumRows() {
		res := *s
		s.Bytes = 0
		s.Record = nil
		return &res
	}

	res := SlicedRecord{
		Record:      s.NewSlice(0, grabRows),
		Bytes:       grabRows * s.bytesPerRow,
		bytesPerRow: s.bytesPerRow,
	}
	s.Record = s.NewSlice(grabRows, s.NumRows()+1)
	s.Bytes -= res.Bytes
	return &res
}

func (s *SlicedRecord) getToFlush(limit *Cap) []arrow.Record {
	// as s.Record != nil we know that the limits are there in place & the s.Record.NumRows() > 0
	grabBytesRows := int64(-1)
	if limit.bytes.limit > 0 {
		grabBytesRows = limit.bytes.limit / s.bytesPerRow
	}
	grabRows := int64(-1)
	if limit.rows.limit > 0 {
		grabRows = limit.rows.limit
	}
	grabRows = min(max(grabRows, 0), max(grabBytesRows, 0))
	if grabRows == 0 {
		// not even a single row fits
		return s.slice()
	}
	if grabRows > s.NumRows() {
		// no need to flush anything, as the amount of rows isn't enough to grant this
		return nil
	}

	flush := make([]arrow.Record, 0, s.NumRows()/grabRows)
	offset := int64(0)
	for offset+grabRows <= s.NumRows() {
		flush = append(flush, s.NewSlice(offset, offset+grabRows))
		offset += grabRows
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
func SliceRecord(r arrow.Record, limit *Cap) (add *SlicedRecord, flush []arrow.Record, remaining *SlicedRecord) {
	return newSlicedRecord(r).split(limit)
}
