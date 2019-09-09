// Copyright © 2019 Alexandre Kovac <contact@kovacou.fr>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package builder

const (
	onKeyword = " ON "
)

// ParseOn create a new Having based on string input.
// This function should be called to initiate the On field from Join.
func ParseOn(str string, args ...interface{}) (out On) {
	out.And(str, args...)
	return
}

// MakeOn create a new On with complex rules.
// This function should be called to initiate the On field from Joins.
func MakeOn(f func(o *On)) (out On) {
	f(&out)
	return
}

// On is the representation of the On clause.
type On struct {
	Where
}
