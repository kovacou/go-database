// Copyright © 2019 Alexandre Kovac <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package builder

import (
	"fmt"
	"strings"
)

// NewColumns create a new Columns.
func NewColumns() Columns {
	return &columns{}
}

// ParseColumns create a new Columns based on a string(s) input.
// This function should be called to initiate the Columns field.
func ParseColumns(cols ...string) Columns {
	out := &columns{}
	return out.Add(cols...)
}

// Columns is list of columns.
type Columns interface {
	// Add a columns to the list.
	Add(...string) Columns

	// String convert Columns to string.
	String() string

	// Len says the size of the string.
	Len() int
}

type columns struct {
	str strings.Builder
}

func (c *columns) Add(col ...string) Columns {
	if len(col) > 0 {
		if c.str.Len() > 0 {
			fmt.Fprintf(&c.str, ",%s", strings.Join(col, ","))
		} else {
			fmt.Fprint(&c.str, strings.Join(col, ","))
		}
	}
	return c
}

func (c *columns) Len() int {
	return c.str.Len()
}

func (c *columns) String() string {
	if c.str.Len() == 0 {
		return "*"
	}

	return c.str.String()
}
