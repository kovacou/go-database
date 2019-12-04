// Copyright © 2019 Alexandre Kovac <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	i := Insert{
		Table: "test",
		Values: H{
			"col1": "val1",
			"col3": "val3",
			"col2": "val2",
		},
	}

	assert.Equal(t, "INSERT INTO test(col1,col2,col3) VALUES(?,?,?)", i.String())
	assert.Equal(t, i.Args(), []interface{}{"val1", "val2", "val3"})

	i.IgnoreMode = true

	assert.Equal(t, "INSERT IGNORE INTO test(col1,col2,col3) VALUES(?,?,?)", i.String())
	assert.Equal(t, i.Args(), []interface{}{"val1", "val2", "val3"})

	i.OnUpdateKeys = Keys{"col1", "col3"}

	assert.Equal(t, "INSERT IGNORE INTO test(col1,col2,col3) VALUES(?,?,?) ON DUPLICATE KEY UPDATE col1 = VALUES(col1),col3 = VALUES(col3)", i.String())
}
