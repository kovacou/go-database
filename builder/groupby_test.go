// Copyright © 2019 Alexandre Kovac <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroupBy(t *testing.T) {
	g := GroupBy{}
	assert.Empty(t, g.String())
	g.Add("col1")
	assert.Equal(t, "col1", g.String())
	g = ParseGroupBy("col2", "col3", "col4")
	assert.Equal(t, "col2,col3,col4", g.String())
}
