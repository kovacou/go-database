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
	whereKeyword = " WHERE "
	andKeyword   = " AND "
	orKeyword    = " OR "
	inKeyword    = " IN "
	notInKeyword = " NOT IN "
)

// ParseWhere create a new Where based on string input.
// This function should be called to initiate the Where field.
func ParseWhere(str string, args ...interface{}) (out Where) {
	out.And(str, args...)
	return
}

// MakeWhere create a new Where with complex rules.
// This function should be called to initiate the Where field.
func MakeWhere(f func(w *Where)) (out Where) {
	f(&out)
	return
}

// Where clause for the SQL query.
type Where struct {
	str  strings.Builder
	args []interface{}
}

// Args return the arguments of the where.
func (w *Where) Args() []interface{} {
	return w.args
}

// String convert Where to string.
func (w *Where) String() string {
	return w.str.String()
}

// And add a new condition "AND".
func (w *Where) And(str string, args ...interface{}) *Where {
	if w.str.Len() > 0 {
		w.str.WriteString(andKeyword)
	}

	w.str.WriteString(str)
	w.args = append(w.args, args...)
	return w
}

// AndIn add a new condition "AND" with the operator IN.
func (w *Where) AndIn(col string, s Slicer) *Where {
	if n := s.Len(); n > 0 {
		w.And(whereIn(col, inKeyword, s), s.S()...)
	}
	return w
}

// AndNotIn add a new condition "AND" with the operator NOT IN.
func (w *Where) AndNotIn(col string, s Slicer) *Where {
	if n := s.Len(); n > 0 {
		w.And(whereIn(col, notInKeyword, s), s.S()...)
	}
	return w
}

// AndWhere merge Where's inside parenthesis with AND condition.
// AND ( where )
func (w *Where) AndWhere(in ...Where) *Where {
	for _, v := range in {
		if v.str.Len() > 0 {
			w.And(fmt.Sprintf("(%s)", v.str.String()), v.args...)
		}
	}
	return w
}

// Or add a new condition "OR".
func (w *Where) Or(str string, args ...interface{}) *Where {
	if w.str.Len() > 0 {
		w.str.WriteString(orKeyword)
	}

	w.str.WriteString(str)
	w.args = append(w.args, args...)
	return w
}

// OrIn add a new condition "OR" with the operator IN.
func (w *Where) OrIn(col string, s Slicer) *Where {
	if n := s.Len(); n > 0 {
		w.Or(whereIn(col, inKeyword, s), s.S()...)
	}
	return w
}

// OrNotIn add a new condition "OR" with the operator NOT IN.
func (w *Where) OrNotIn(col string, s Slicer) *Where {
	if n := s.Len(); n > 0 {
		w.Or(whereIn(col, notInKeyword, s), s.S()...)
	}
	return w
}

// OrWhere merge Where's inside parenthesis with OR condition.
// OR ( where )
func (w *Where) OrWhere(in ...Where) *Where {
	for _, v := range in {
		if v.str.Len() > 0 {
			w.Or(fmt.Sprintf("(%s)", v.str.String()), v.args...)
		}
	}
	return w
}

// whereIn
func whereIn(col, keyword string, s Slicer) string {
	return fmt.Sprintf("%s%s(%s)", col, keyword, strings.TrimRight(strings.Repeat("?,", s.Len()), ","))
}
