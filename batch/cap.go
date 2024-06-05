package batch

type capped struct {
	current, limit int64
}

func (c capped) reachedLimit() bool { return c.limit > 0 && c.current >= c.limit }

type Cap struct {
	bytes, rows capped
}

func (c *Cap) ReachedLimit() bool { return c.bytes.reachedLimit() || c.rows.reachedLimit() }

func (c *Cap) Rows() int64 { return c.rows.current }

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
