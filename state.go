// Copyright © 2019 Alexandre Kovac <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package database

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// QueryState is the interface that abstract the state of the execution of an query.
type QueryState interface {
	ContextID() string
	ContextFlag() []string
	Start() time.Time
	End() time.Time
	Runtime() time.Duration
	String() string
	Bytes() []byte
}

type qs struct {
	query   string
	args    []interface{}
	ctxID   string
	ctxFlag []string
	start   time.Time
	end     time.Time
}

// Runtime
func (p *qs) Runtime() time.Duration {
	return p.end.Sub(p.start)
}

// Start
func (p *qs) Start() time.Time {
	return p.start
}

// End
func (p *qs) End() time.Time {
	return p.end
}

// ContextID
func (p *qs) ContextID() string {
	return p.ctxID
}

// ContextFlag
func (p *qs) ContextFlag() []string {
	return p.ctxFlag
}

// Bytes
func (p *qs) Bytes() []byte {
	b := bytes.NewBufferString("")
	q := p.query

	for _, a := range p.args {
		arg := fmt.Sprint(a)

		if _, err := strconv.ParseInt(arg, 0, 0); err != nil {
			if _, err := strconv.ParseFloat(arg, 64); err != nil {
				arg = fmt.Sprintf("\"%s\"", arg)
			}
		}

		q = strings.Replace(q, "?", arg, 1)
	}

	b.WriteString(q)
	return b.Bytes()
}

// String
func (p *qs) String() string {
	return string(p.Bytes())
}
