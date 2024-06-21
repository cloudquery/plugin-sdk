package batch

type capped struct {
	current, limit int64
}

func (c capped) reachedLimit() bool { return c.limit > 0 && c.current >= c.limit }
func (c capped) remaining() int64 {
	if c.limit > 0 {
		return c.limit - c.current
	}
	return -1
}

func (c capped) remainingPerN(n int64) int64 {
	if c.limit > 0 {
		return (c.limit - c.current) / n
	}
	return -1
}

func (c capped) cap() int64 {
	if c.limit > 0 {
		return c.limit
	}
	return -1
}
func (c capped) capPerN(n int64) int64 {
	if c.limit > 0 {
		return c.limit / n
	}
	return -1
}

type Cap struct {
	bytes, rows capped
}

func (c *Cap) ReachedLimit() bool { return c.bytes.reachedLimit() || c.rows.reachedLimit() }
func (c *Cap) Rows() int64        { return c.rows.current }
func (c *Cap) AddRows(rows int64) { c.rows.current += rows }

func (c *Cap) AddSlice(record *SlicedRecord) {
	c.rows.current += record.NumRows()
	c.bytes.current += record.Bytes
}

func (c *Cap) Reset() {
	c.bytes.current = 0
	c.rows.current = 0
}

func (c *Cap) add(bytes, rows int64) {
	c.bytes.current += bytes
	c.rows.current += rows
}

func (c *Cap) set(bytes, rows int64) {
	c.bytes.current = bytes
	c.rows.current = rows
}

func CappedAt(bytes, rows int64) *Cap {
	return &Cap{
		bytes: capped{limit: bytes},
		rows:  capped{limit: rows},
	}
}
