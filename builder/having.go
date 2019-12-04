// Copyright © 2019 Alexandre Kovac <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package builder

var (
	havingKeyword = " HAVING "
)

// ParseHaving create a new Having based on string input.
// This function should be called to initiate the Having field.
func ParseHaving(str string, args ...interface{}) (out Having) {
	out.And(str, args...)
	return
}

// MakeHaving create a new Having with complex rules.
// This function should be called to initiate the Having field.
func MakeHaving(f func(o *Having)) (out Having) {
	f(&out)
	return
}

// Having is the representation of the Having clause.
type Having struct {
	Where
}
