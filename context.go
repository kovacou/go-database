// Copyright © 2019 Alexandre Kovac <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package database

import (
	"strings"

	"github.com/rs/xid"
)

// ContextFunc handler.
type ContextFunc func(Connection) error

// Context is the interface that abstract all interaction of the connection.
// Mainly used to bypass the unexported ctx.
type Context interface {
	ID() string
	Done()
	Flag() string
	Len() int
	Flush()
	Err() error
	Push(QueryState) bool
}

// newContext create a new context.
func newContext(f []string) *ctx {
	ctx := &ctx{
		flag:  f,
		store: make(chan QueryState, 1),
		done:  make(chan struct{}),
		qsl:   make([]QueryState, 0),
	}

	ctx.generateID()
	go func() {
		for {
			select {
			case qs := <-ctx.store:
				ctx.qsl = append(ctx.qsl, qs)
			case <-ctx.done:
				close(ctx.store)
				close(ctx.done)
				return
			}
		}
	}()

	return ctx
}

type ctx struct {
	id    string
	flag  []string
	store chan QueryState
	done  chan struct{}
	qsl   []QueryState
}

// generateID generate a new unique id for the context.
func (ctx *ctx) generateID() {
	ctx.id = xid.New().String()
}

// ID returns the id of the context.
func (ctx *ctx) ID() string {
	return ctx.id
}

// Flag will compute the internal flag.
func (ctx *ctx) Flag() string {
	return strings.Join(ctx.flag, "")
}

// Done the context.
func (ctx *ctx) Done() {
	if ctx.done != nil {
		ctx.done <- struct{}{}
	}
}

// Flush clear the list of query state.
func (ctx *ctx) Flush() {
	ctx.qsl = []QueryState{}
}

// Len returns the number of query state into the context.
func (ctx *ctx) Len() int {
	return len(ctx.qsl)
}

// Err return the error associated to the context.
func (ctx *ctx) Err() error {
	return nil
}

// Push a new QueryState into the context.
func (ctx *ctx) Push(qs QueryState) bool {
	if ctx.store != nil {
		ctx.store <- qs
		return true
	}
	return false
}

// -------------------------------------------------

// HasContext says if the connection has a context.
func (conn *db) HasContext() bool {
	return conn.ctx != nil
}

// Context copy the current connection and asign a new Context.
func (conn *db) Context(f ...string) Connection {
	connCtx := conn.copy()
	connCtx.ctx = newContext(f)
	return connCtx
}

// RunContext run a bunch of ContextFunc and handle the error.
func (conn *db) RunContext(funcs ...ContextFunc) (err error) {
	ctx := conn.Context()
	for _, f := range funcs {
		if err = f(ctx); err != nil {
			break
		}
	}
	ctx.Done()
	return
}

// Done terminate the context.
func (conn *db) Done() {
	if conn.HasContext() {
		conn.ctx.Done()

		if conn.hasProfiling() {
			// 	conn.Profile()
		}
	}
}
