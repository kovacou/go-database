// Copyright © 2019 Alexandre Kovac <contact@kovacou.fr>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHKeys(t *testing.T) {
	h := H{
		"key1": "v1",
		"key3": "v3",
		"key2": "v2",
		"key0": "v0",
	}

	assert.Equal(t, h.Keys(), Keys{"key0", "key1", "key2", "key3"})
}

func TestBindsString(t *testing.T) {
	b := Binds{
		"col1": "col1 + 15",
	}

	assert.Equal(t, b.String(), "col1 = col1 + 15")
}
