// Copyright © 2019 Alexandre Kovac <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package builder

const (
	onKeyword = " ON "
)

// NewOn create a new On.
func NewOn() On {
	return &where{}
}

// ParseOn create a new Having based on string input.
// This function should be called to initiate the On field from Join.
func ParseOn(str string, args ...any) On {
	out := &where{}
	return out.And(str, args...)
}

// MakeOn create a new On with complex rules.
// This function should be called to initiate the On field from Joins.
func MakeOn(f func(o On)) On {
	out := &where{}
	f(out)
	return out
}

// On is the representation of the On clause.
type On interface {
	Where
}
