// Copyright © 2019 Alexandre Kovac <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package builder

import (
	"strings"
)

// NewQuery create a new Query based on string input.
func NewQuery(str string, args []any) *Query {
	q := &Query{
		args: args,
	}
	q.str.WriteString(str)
	return q
}

// Query is the representation of an Query statement.
type Query struct {
	str  strings.Builder
	args []any
}

// String convert query to string.
func (q *Query) String() string {
	return q.str.String()
}

// Args return the arguments of the query.
func (q *Query) Args() []any {
	return q.args
}
