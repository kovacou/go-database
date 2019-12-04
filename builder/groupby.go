// Copyright © 2019 Alexandre Kovac <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package builder

const (
	groupByKeyword = " GROUP BY "
)

// ParseGroupBy create a new GroupBy based on a string(s) input.
// This function should be called to initiate the GroupBy field.
func ParseGroupBy(cols ...string) (out GroupBy) {
	out.Add(cols...)
	return
}

// GroupBy is the representation of the GROUP BY clause.
type GroupBy struct {
	Columns
}

// String return the natural string of Columns with "*".
func (gb *GroupBy) String() string {
	return gb.str.String()
}
