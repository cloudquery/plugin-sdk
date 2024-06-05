package batch

import "golang.org/x/exp/constraints"

type Cap[T constraints.Integer | constraints.Float] struct {
	current, limit T
}

func (c *Cap[T]) Add(v T) (reachedLimit bool) {
	if c != nil {
		c.current += v
	}
	return c.ReachedLimit()
}

func (c *Cap[T]) ReachedLimit() bool {
	if c == nil {
		return false
	}
	return c.limit > 0 && c.current >= c.limit
}
