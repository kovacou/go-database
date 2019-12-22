// Copyright © 2019 Alexandre Kovac <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeWhere(t *testing.T) {
	w := MakeWhere(func(w Where) {
		w.And("col1 = ?", "test")
		w.And("col2 IS NOT NULL")
	})

	assert.Equal(t, w.String(), "col1 = ? AND col2 IS NOT NULL")
	assert.Len(t, w.Args(), 1)
	assert.Contains(t, w.Args(), "test")
}

func TestParseWhere(t *testing.T) {
	w := ParseWhere("col1 = ? OR col2", "value2")

	assert.Equal(t, "col1 = ? OR col2", w.String())
	assert.Len(t, w.Args(), 1)
}
func TestWhereAnd(t *testing.T) {
	w := where{}
	w.And("col1 = ?", "test")
	assert.Equal(t, "col1 = ?", w.String())
	assert.Equal(t, []interface{}{"test"}, w.Args())
	w.And("col2")
	assert.Equal(t, "col1 = ? AND col2", w.String())
	assert.Equal(t, []interface{}{"test"}, w.Args())
}

func TestWhereOr(t *testing.T) {
	w := where{}
	w.Or("col1 = ?", "test")
	assert.Equal(t, "col1 = ?", w.String())
	assert.Equal(t, []interface{}{"test"}, w.Args())
	w.Or("col2 = ?", "test2")
	assert.Equal(t, "col1 = ? OR col2 = ?", w.String())
	assert.Equal(t, []interface{}{"test", "test2"}, w.Args())
}

func TestWhereAndOr(t *testing.T) {
	w := where{}
	w.And("col1 = ?")
	assert.Equal(t, "col1 = ?", w.String())
	w.Or("col2")
	assert.Equal(t, "col1 = ? OR col2", w.String())
}
