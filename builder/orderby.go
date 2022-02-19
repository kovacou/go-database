// Copyright © 2019 Alexandre Kovac <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package builder

const (
	orderByKeyword = " ORDER BY "
)

// NewOrderBy create a new OrderBy.
func NewOrderBy() OrderBy {
	return &orderBy{}
}

// ParseOrderBy create a new OrderBy based on a string(s) input.
// This function should be called to initiate the OrderBy field.
func ParseOrderBy(cols ...string) OrderBy {
	out := &orderBy{}
	return out.Add(cols...)
}

// ASC returns ASC if v is true else DESC.
func ASC(v bool) string {
	if v {
		return "ASC"
	}
	return "DESC"
}

// DESC returns DESC if v is true else ASC.
func DESC(v bool) string {
	if v {
		return "DESC"
	}
	return "ASC"
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
