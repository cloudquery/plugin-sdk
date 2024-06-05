package batch

type Cap struct {
	current, limit int64
}

func (c *Cap) ReachedLimit() bool { return c.limit > 0 && c.current >= c.limit }
func (c *Cap) Reset()             { c.current = 0 }
func (c *Cap) Current() int64     { return c.current }

func CappedAt(limit int64) *Cap { return &Cap{limit: limit} }
