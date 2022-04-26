// Copyright © 2019 Alexandre Kovac <contact@kovacou.com>.
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
	Limit   uint64
	Offset  uint64
}

// String convert the select to string.
func (s *Select) String() string {
	q := strings.Builder{}
	q.WriteString(selectKeyword)
	if s.Columns != nil {
		q.WriteString(s.Columns.String())
	} else {
		q.WriteString("*")
	}
	q.WriteString(fromKeyword)
	q.WriteString(s.Table)

	// Joins section
	if s.Joins.Len() > 0 {
		q.WriteString(s.Joins.String())
	}

	// Where clause
	if s.Where != nil && s.Where.Len() > 0 {
		q.WriteString(whereKeyword)
		q.WriteString(s.Where.String())
	}

	// Group By clause
	if s.GroupBy != nil && s.GroupBy.Len() > 0 {
		q.WriteString(groupByKeyword)
		q.WriteString(s.GroupBy.String())
	}

	// Having clause
	if s.Having != nil && s.Having.Len() > 0 {
		q.WriteString(havingKeyword)
		q.WriteString(s.Having.String())
	}

	// OrderBy clause
	if s.OrderBy != nil && s.OrderBy.Len() > 0 {
		q.WriteString(orderByKeyword)
		q.WriteString(s.OrderBy.String())
	}

	// Pagination
	if s.Limit > 0 {
		fmt.Fprintf(&q, " LIMIT %d OFFSET %d ", s.Limit, s.Offset)
	}

	return q.String()
}

// Args compute the arguments of the select query.
func (s *Select) Args() (out []any) {
	out = append(out, s.Joins.Args()...)
	if s.Where != nil {
		out = append(out, s.Where.Args()...)
	}
	if s.Having != nil {
		out = append(out, s.Having.Args()...)
	}
	return
}
