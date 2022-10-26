// Copyright © 2019 Alexandre Kovac <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package builder

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	whereKeyword = " WHERE "
	andKeyword   = " AND "
	orKeyword    = " OR "
	inKeyword    = " IN "
	notInKeyword = " NOT IN "
)

// NewWhere create a new Where.
func NewWhere() Where {
	return &where{}
}

// ParseWhere create a new Where based on string input.
// This function should be called to initiate the Where field.
func ParseWhere(str string, args ...any) Where {
	out := &where{}
	return out.And(str, args...)
}

// MakeWhere create a new Where with complex rules.
// This function should be called to initiate the Where field.
func MakeWhere(f func(w Where)) Where {
	out := &where{}
	f(out)
	return out
}

// Where clause for the SQL query.
type Where interface {
	// Args return the arguments of the where.
	Args() []any

	// String convert Where to string.
	String() string

	// And add a new condition "AND".
	And(str string, args ...any) Where

	// AndIf a new condition "AND" if args is Valid & and not Zero.
	AndIf(str string, arg any) Where

	// AndIn add a new condition "AND" with the operator IN.
	AndIn(col string, s Slicer) Where

	// AndNotIn add a new condition "AND" with the operator NOT IN.
	AndNotIn(col string, s Slicer) Where

	// AndWhere merge Where's inside parenthesis with AND condition.
	// AND ( where )
	AndWhere(in ...Where) Where

	// Or add a new condition "OR".
	Or(str string, args ...any) Where

	// OrIf a new condition "or" if args is Valid & and not Zero.
	OrIf(str string, arg any) Where

	// OrIn add a new condition "OR" with the operator IN.
	OrIn(col string, s Slicer) Where

	// OrNotIn add a new condition "OR" with the operator NOT IN.
	OrNotIn(col string, s Slicer) Where

	// OrWhere merge Where's inside parenthesis with OR condition.
	// OR ( where )
	OrWhere(in ...Where) Where

	// Len says the length of the string.
	Len() int
}

type where struct {
	str  strings.Builder
	args []any
}

func (w *where) Args() []any {
	return w.args
}

func (w *where) String() string {
	return w.str.String()
}

func (w *where) And(str string, args ...any) Where {
	if w.str.Len() > 0 {
		w.str.WriteString(andKeyword)
	}

	w.str.WriteString(str)
	w.args = append(w.args, args...)
	return w
}

func (w *where) AndIf(str string, arg any) Where {
	if !reflect.ValueOf(arg).IsZero() {
		w.And(str, arg)
	}
	return w
}

func (w *where) AndIn(col string, s Slicer) Where {
	if n := s.Len(); n > 0 {
		w.And(whereIn(col, inKeyword, s), s.S()...)
	}
	return w
}

func (w *where) AndNotIn(col string, s Slicer) Where {
	if n := s.Len(); n > 0 {
		w.And(whereIn(col, notInKeyword, s), s.S()...)
	}
	return w
}

func (w *where) AndWhere(in ...Where) Where {
	for _, v := range in {
		if v.Len() > 0 {
			w.And(fmt.Sprintf("(%s)", v.String()), v.Args()...)
		}
	}
	return w
}

func (w *where) Or(str string, args ...any) Where {
	if w.str.Len() > 0 {
		w.str.WriteString(orKeyword)
	}

	w.str.WriteString(str)
	w.args = append(w.args, args...)
	return w
}

func (w *where) OrIf(str string, arg any) Where {
	if !reflect.ValueOf(arg).IsZero() {
		w.Or(str, arg)
	}

	return w
}

func (w *where) OrIn(col string, s Slicer) Where {
	if n := s.Len(); n > 0 {
		w.Or(whereIn(col, inKeyword, s), s.S()...)
	}
	return w
}

func (w *where) OrNotIn(col string, s Slicer) Where {
	if n := s.Len(); n > 0 {
		w.Or(whereIn(col, notInKeyword, s), s.S()...)
	}
	return w
}

func (w *where) OrWhere(in ...Where) Where {
	for _, v := range in {
		if v.Len() > 0 {
			w.Or(fmt.Sprintf("(%s)", v.String()), v.Args()...)
		}
	}
	return w
}

func (w *where) Len() int {
	return w.str.Len()
}

// whereIn
func whereIn(col, keyword string, s Slicer) string {
	return fmt.Sprintf("%s%s(%s)", col, keyword, strings.TrimRight(strings.Repeat("?,", s.Len()), ","))
}
