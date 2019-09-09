// Copyright © 2019 Alexandre Kovac <contact@kovacou.fr>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package database

import "time"

// QueryState is the interface that abstract the state of the execution of an query.
type QueryState interface {
	Start() time.Time
	End() time.Time
	Runtime() time.Duration
	String() string
	Bytes() []byte
}
