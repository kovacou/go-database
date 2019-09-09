// Copyright © 2019 Alexandre Kovac <contact@kovacou.fr>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColumns(t *testing.T) {
	cols := Columns{}
	assert.Equal(t, "*", cols.String())
	cols.Add("col1")
	assert.Equal(t, "col1", cols.String())
	cols.Add("col2")
	assert.Equal(t, "col1,col2", cols.String())
}
