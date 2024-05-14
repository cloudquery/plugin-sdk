package writers

import "golang.org/x/exp/constraints"

type Capped[A constraints.Integer] struct {
	curr, limit A
}

func (c *Capped[A]) Add(amount A)              { c.curr += amount }
func (c *Capped[A]) Reset()                    { c.curr = 0 }
func (c *Capped[A]) Current() A                { return c.curr }
func (c *Capped[A]) OverflownBy(amount A) bool { return c.limit > 0 && c.curr+amount > c.limit }

func NewCapped[A constraints.Integer](limit A) Capped[A] {
	return Capped[A]{limit: limit}
}
