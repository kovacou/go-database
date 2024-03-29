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
	updateKeyword = "UPDATE "
	setKeyword    = " SET "
)

// NewUpdate create a new update.
func NewUpdate(t string) *Update {
	return &Update{
		Table:  t,
		Values: H{},
		Binds:  Binds{},
		Joins:  Joins{},
		Where:  NewWhere(),
	}
}

// Update is the representation of an Update statement.
type Update struct {
	Table  string
	Values H
	Binds  Binds
	Joins  Joins
	Where  Where
}

// String convert the update to string.
func (u *Update) String() string {
	q := strings.Builder{}
	q.WriteString(updateKeyword)
	q.WriteString(u.Table)
	q.WriteString(setKeyword)

	// Values
	{
		keys := u.Values.Keys()
		n := len(keys) - 1

		for i, k := range keys {
			if i < n {
				fmt.Fprintf(&q, "%s = ?,", k)
			} else {
				fmt.Fprintf(&q, "%s = ?", k)
			}
		}
	}

	// Joins section
	if u.Joins.Len() > 0 {
		q.WriteString(u.Joins.String())
	}

	// WHERE clause
	if u.Where != nil && u.Where.Len() > 0 {
		q.WriteString(whereKeyword)
		q.WriteString(u.Where.String())
	}
	return q.String()
}

// Args compute the arguments of the update statement.
func (u *Update) Args() (out []any) {
	for _, k := range u.Values.Keys() {
		out = append(out, u.Values[k])
	}

	if u.Where != nil {
		out = append(out, u.Where.Args()...)
	}
	return
}
