// Copyright © 2019 Alexandre Kovac <contact@kovacou.fr>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package database

// Profiler is an interface that abstract all interactions with the profiler.
type Profiler interface {
	Close() error
}

type profiler struct {
}

// -------------------------------------------------

// HasProfiler says if the connection has a profiler.
func (conn *db) HasProfiler() bool {
	return conn.profiler != nil
}
