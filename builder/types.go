// Copyright © 2019 Alexandre Kovac <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package builder

import (
	"fmt"
	"sort"
	"strings"
)

// Slicer abstract a slice (until generics are supported)
type Slicer interface {
	Len() int
	S() []interface{}
}

// Keys is a slice of string representing a list of key.
type Keys []string

// RawKeys is used to map keys to expression (On Duplicate Key Update case).
type RawKeys map[string]string

// H is the representation of an insert or update values.
type H map[string]interface{}

// Keys return sorted keys of the values.
func (v H) Keys() Keys {
	keys := make(Keys, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Binds is used to bind column to another one. (Raw version, no secure)
type Binds map[string]string

// String
func (b Binds) String() string {
	str := strings.Builder{}
	for k, v := range b {
		fmt.Fprintf(&str, "%s = %s", k, v)
	}
	return str.String()
}
