// Copyright © 2019 Alexandre Kovac <contact@kovacou.fr>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package builder

const (
	groupByKeyword = " GROUP BY "
)

// GroupBy is the representation of the GROUP BY clause.
type GroupBy struct {
	Columns
}

// String return the natural string of Columns with "*".
func (gb *GroupBy) String() string {
	return gb.str.String()
}
