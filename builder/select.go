// Copyright © 2019 Alexandre Kovac <contact@kovacou.fr>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package builder

import (
	"fmt"
	"strings"
)

const (
	selectKeyword = " SELECT "
	fromKeyword   = " FROM "
)

// Select is the representation of the Select statement.
type Select struct {
	Table   string
	Columns Columns
	Joins   Joins
	Where   Where
	Having  Having
	GroupBy GroupBy
	OrderBy OrderBy
	Limit   int
	Offset  int
}

// String convert the select to string.
func (s *Select) String() string {
	q := strings.Builder{}
	q.WriteString(selectKeyword)
	q.WriteString(s.Columns.String())
	q.WriteString(fromKeyword)
	q.WriteString(s.Table)

	// Joins section
	if s.Joins.Len() > 0 {
		q.WriteString(s.Joins.String())
	}

	// Where clause
	if s.Where.str.Len() > 0 {
		q.WriteString(whereKeyword)
		q.WriteString(s.Where.str.String())
	}

	// Group By clause
	if s.GroupBy.str.Len() > 0 {
		q.WriteString(groupByKeyword)
		q.WriteString(s.GroupBy.str.String())
	}

	// Having clause
	if s.Having.str.Len() > 0 {
		q.WriteString(havingKeyword)
		q.WriteString(s.Having.str.String())
	}

	// OrderBy clause
	if s.OrderBy.str.Len() > 0 {
		q.WriteString(orderByKeyword)
		q.WriteString(s.OrderBy.str.String())
	}

	// Pagination
	if s.Limit > 0 {
		fmt.Fprintf(&q, " LIMIT %d OFFSET %d ", s.Limit, s.Offset)
	}

	return q.String()
}

// Args compute the arguments of the select query.
func (s *Select) Args() (out []interface{}) {
	out = append(out, s.Joins.Args()...)
	out = append(out, s.Where.args...)
	out = append(out, s.Having.args...)
	return
}
