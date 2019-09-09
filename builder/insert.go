// Copyright © 2019 Alexandre Kovac <contact@kovacou.fr>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package builder

import (
	"fmt"
	"sort"
	"strings"
)

const (
	insertKeyword               = "INSERT "
	ignoreKeyword               = "IGNORE "
	intoKeyword                 = "INTO "
	valuesKeyword               = " VALUES"
	onDuplicateKeyUpdateKeyword = " ON DUPLICATE KEY UPDATE "
)

// Values is the representation of an insert or update values.
type Values map[string]interface{}

// Keys is a slice of string representing a list of key.
type Keys []string

// RawKeys is used to map keys to expression (On Duplicate Key Update case).
type RawKeys map[string]string

// Keys return sorted keys of the values.
func (v Values) Keys() Keys {
	keys := make(Keys, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Insert is the representation of an Insert statement.
type Insert struct {
	Table           string
	IgnoreMode      bool
	Values          Values
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

	q.WriteString(intoKeyword)
	q.WriteString(i.Table)
	q.WriteRune('(')
	keys := i.Values.Keys()
	n := len(keys)
	q.WriteString(strings.Join(keys, ","))
	q.WriteRune(')')
	q.WriteString(valuesKeyword)
	q.WriteRune('(')
	q.WriteString(strings.Repeat("?,", n)[:(n*2)-1])
	q.WriteRune(')')

	if len(i.OnUpdateKeys) > 0 || len(i.OnUpdateRawKeys) > 0 {
		q.WriteString(onDuplicateKeyUpdateKeyword)
		for _, k := range i.OnUpdateKeys {
			fmt.Fprintf(&q, "%s = VALUES(%s)", k, k)
		}

		if i.OnUpdateRawKeys != nil {
			for key, exp := range i.OnUpdateRawKeys {
				fmt.Fprintf(&q, "%s = %s", key, exp)
			}
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
