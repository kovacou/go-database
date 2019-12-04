// Copyright © 2019 Alexandre Kovac <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	d := Delete{
		Table: "test",
		Where: ParseWhere("col1 = ? AND col2 = ?", 1, "val"),
	}

	assert.Equal(t, "DELETE  FROM test WHERE col1 = ? AND col2 = ?", d.String())
	assert.Equal(t, d.Args(), []interface{}{1, "val"})
}
