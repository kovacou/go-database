// Copyright © 2019 Alexandre Kovac <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package builder

import (
	"strings"
)

const (
	deleteKeyword = "DELETE "
)

// Delete is the representation of an Delete statement.
type Delete struct {
	Table string
	Where Where
}

// String convert the delete into string.
func (d *Delete) String() string {
	q := strings.Builder{}
	q.WriteString(deleteKeyword)
	q.WriteString(fromKeyword)
	q.WriteString(d.Table)

	if d.Where != nil && d.Where.Len() > 0 {
		q.WriteString(whereKeyword)
		q.WriteString(d.Where.String())
	}
	return q.String()
}

// Args compute the arguments of the delete statement.
func (d *Delete) Args() (out []interface{}) {
	if d.Where != nil {
		out = d.Where.Args()
	}
	return
}
