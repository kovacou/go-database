// Copyright © 2019 Alexandre Kovac <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package builder

import (
	"fmt"
	"strings"
)

// ParseColumns create a new Columns based on a string(s) input.
// This function should be called to initiate the Columns field.
func ParseColumns(cols ...string) (out Columns) {
	out.Add(cols...)
	return
}

// Columns is list of columns.
type Columns struct {
	str strings.Builder
}

// Add a columns to the list.
func (c *Columns) Add(col ...string) *Columns {
	if len(col) > 0 {
		if c.str.Len() > 0 {
			fmt.Fprintf(&c.str, ",%s", strings.Join(col, ","))
		} else {
			fmt.Fprint(&c.str, strings.Join(col, ","))
		}
	}
	return c
}

// String convert Columns to string.
func (c *Columns) String() string {
	if c.str.Len() == 0 {
		return "*"
	}

	return c.str.String()
}
