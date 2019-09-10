// Copyright © 2019 Alexandre Kovac <contact@kovacou.fr>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package builder

const (
	orderByKeyword = " ORDER BY "
)

// ParseOrderBy create a new OrderBy based on a string(s) input.
// This function should be called to initiate the OrderBy field.
func ParseOrderBy(cols ...string) (out OrderBy) {
	out.Add(cols...)
	return
}

// OrderBy clause for the SELECT query.
type OrderBy struct {
	Columns
}

// String return the natural string of Columns with "*".
func (gb *OrderBy) String() string {
	return gb.str.String()
}
