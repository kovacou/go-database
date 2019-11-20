// Copyright © 2019 Alexandre Kovac <contact@kovacou.fr>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package database

import (
	"strconv"
	"time"
)

// toBool convert string v to bool.
func toBool(v string) bool {
	out, err := strconv.ParseBool(v)
	if err != nil && len(v) > 0 {
		out = true
	}
	return out
}

// toInt convert string v to int.
func toInt(v string) int {
	out, _ := strconv.ParseInt(v, 0, 0)
	return int(out)
}

// toDuration convert string to duration.
func toDuration(v string) time.Duration {
	out, _ := time.ParseDuration(v)
	return out
}
