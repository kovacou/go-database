// Copyright © 2019 Alexandre Kovac <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package builder

const (
	orderByKeyword = " ORDER BY "
)

// ParseOrderBy create a new OrderBy based on a string(s) input.
// This function should be called to initiate the OrderBy field.
func ParseOrderBy(cols ...string) OrderBy {
	out := &orderBy{}
	return out.Add(cols...)
}

// OrderBy clause for the SELECT query.
type OrderBy interface {
	Columns
}

type orderBy struct {
	columns
}

func (ob *orderBy) String() string {
	return ob.str.String()
}
