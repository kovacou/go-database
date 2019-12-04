// Copyright © 2019 Alexandre Kovac <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package builder

import (
	"fmt"
	"strings"
)

// Joins is a slice of Join.
type Joins []Join

// Add a new Join to the list.
func (j *Joins) Add(joins ...Join) *Joins {
	*j = append(*j, joins...)
	return j
}

// Len return the number of join.
func (j Joins) Len() int {
	return len(j)
}

// String convert Joins to string.
func (j Joins) String() string {
	str := strings.Builder{}
	for i := range j {
		str.WriteString(j[i].String())
	}
	return str.String()
}

// Args return the args of the joins.
func (j Joins) Args() []interface{} {
	out := []interface{}{}
	for i := range j {
		out = append(out, j[i].On.args...)
	}
	return out
}

// Join is the representation of the Join.
type Join struct {
	Table string
	Type  string
	On    On
}

// String convert Join to string.
func (j *Join) String() string {
	str := strings.Builder{}
	if j.Type != "" {
		fmt.Fprintf(&str, " %s JOIN %s", j.Type, j.Table)
	} else {
		fmt.Fprintf(&str, " JOIN %s", j.Table)
	}

	if j.On.str.Len() > 0 {
		str.WriteString(onKeyword)
		str.WriteString(j.On.str.String())
	}
	return str.String()
}

// Args return the arguments for the join.
func (j *Join) Args() []interface{} {
	return j.On.args
}
