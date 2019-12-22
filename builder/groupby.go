// Copyright © 2019 Alexandre Kovac <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package builder

const (
	groupByKeyword = " GROUP BY "
)

// NewGroupBy create a new GroupBy.
func NewGroupBy() GroupBy {
	return &groupBy{}
}

// ParseGroupBy create a new GroupBy based on a string(s) input.
// This function should be called to initiate the GroupBy field.
func ParseGroupBy(cols ...string) GroupBy {
	out := &groupBy{}
	return out.Add(cols...)
}

// GroupBy is the representation of the GROUP BY clause.
type GroupBy interface {
	Columns
}

type groupBy struct {
	columns
}

// String return the natural string of Columns with "*".
func (gb *groupBy) String() string {
	return gb.str.String()
}
