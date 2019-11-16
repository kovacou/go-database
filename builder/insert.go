// Copyright © 2019 Alexandre Kovac <contact@kovacou.fr>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package builder

import (
	"fmt"
	"strings"
)

const (
	insertKeyword               = "INSERT "
	ignoreKeyword               = "IGNORE "
	intoKeyword                 = "INTO "
	valuesKeyword               = " VALUES"
	onDuplicateKeyUpdateKeyword = " ON DUPLICATE KEY UPDATE "
)

// Insert is the representation of an Insert statement.
type Insert struct {
	Table           string
	Select          Select
	IgnoreMode      bool
	Values          H
	keys            Keys
	OnUpdateKeys    Keys
	OnUpdateRawKeys RawKeys
}

// String convert the insert to string.
func (i *Insert) String() string {
	q := strings.Builder{}
	q.WriteString(insertKeyword)

	if i.IgnoreMode {
		q.WriteString(ignoreKeyword)
	}

	keys := i.Values.Keys()
	n := len(keys)

	q.WriteString(intoKeyword)
	q.WriteString(i.Table)
	q.WriteRune('(')
	q.WriteString(strings.Join(keys, ","))
	q.WriteRune(')')
	q.WriteString(valuesKeyword)
	q.WriteRune('(')
	q.WriteString(strings.Repeat("?,", n)[:(n*2)-1])
	q.WriteRune(')')

	if len(i.OnUpdateKeys) > 0 || len(i.OnUpdateRawKeys) > 0 {
		q.WriteString(onDuplicateKeyUpdateKeyword)
		first := true

		for _, k := range i.OnUpdateKeys {
			if !first {
				q.WriteRune(',')
			} else {
				first = false
			}

			fmt.Fprintf(&q, "%s = VALUES(%s)", k, k)
		}

		for key, exp := range i.OnUpdateRawKeys {
			if !first {
				q.WriteRune(',')
			} else {
				first = false
			}

			fmt.Fprintf(&q, "%s = %s", key, exp)
		}
	}
	return q.String()
}

// Args compute the arguments of the insert statement.
func (i *Insert) Args() (out []interface{}) {
	for _, k := range i.Values.Keys() {
		out = append(out, i.Values[k])
	}
	return
}
