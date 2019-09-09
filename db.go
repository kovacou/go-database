// Copyright © 2019 Alexandre Kovac <contact@kovacou.fr>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package database

import (
	"database/sql/driver"
	"log"
	"sync"

	"github.com/jmoiron/sqlx"
)

// db is a wrapper around sqlx.DB.
type db struct {
	id  uint
	dbx **sqlx.DB
	tx  *sqlx.Tx
	m   *sync.Mutex
	log *log.Logger
	err error

	ctx      Context
	env      Environment
	profiler Profiler
}

// Copy the current connection.
func (conn *db) Copy() Connection {
	return conn
}

// copy the current connection.
func (conn *db) copy() *db {
	return conn
}

// SetLogger set a new logger to the db.
func (conn *db) SetLogger(log *log.Logger) {
	if log != nil {
		conn.log = log
	}
}

// hasDebug says if the connection have debug mode enable.
func (conn *db) hasDebug() bool {
	return Debug || conn.env.Debug
}

// hasVerbose says if the connection have verbose mode enable.
func (conn *db) hasVerbose() bool {
	return Verbose || conn.env.Verbose
}

// DB return the unwrapped sqlx.DB.
func (conn *db) DB() *sqlx.DB {
	conn.Connect()
	return *conn.dbx
}

// Close closes the database and prevents new queries from starting.
func (conn *db) Close() (err error) {
	if conn.ctx != nil {
		conn.ctx.Done()
	}

	if conn.profiler != nil {
		conn.profiler.Close()
	}

	if dbx := (*conn.dbx); dbx != nil {
		err = dbx.Close()
	}
	return
}

// Ping verifies a connection to the database is still alive,
// establishing a connection if necessary.
func (conn *db) Ping() error {
	if err := conn.Connect(); err != nil {
		return err
	}

	if dbx := (*conn.dbx); dbx != nil {
		return dbx.Ping()
	}
	return driver.ErrBadConn
}

// MustPing call Ping and panic if there is an error.
func (conn *db) MustPing() {
	if err := conn.Ping(); err != nil {
		panic(err.Error())
	}
}

// Connect to a database and verify with a ping.
func (conn *db) Connect() (err error) {
	if (*conn.dbx) != nil {
		return
	}

	conn.m.Lock()
	if (*conn.dbx) != nil {
		conn.m.Unlock()
		return
	}

	(*conn.dbx), err = sqlx.Connect(
		conn.env.Driver,
		conn.env.String(),
	)

	conn.m.Unlock()
	if err != nil {
		return
	}

	dbx := (*conn.dbx)
	dbx.SetMaxIdleConns(conn.env.MaxIdle)
	dbx.SetMaxOpenConns(conn.env.MaxOpen)
	dbx.SetConnMaxLifetime(conn.env.MaxLifetime)

	if conn.hasVerbose() {
		conn.log.Printf("connected")
	}

	return
}

// LastError
func (conn *db) LastError() error {
	return conn.err
}