// Copyright © 2019 Alexandre Kovac <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/jmoiron/sqlx"
)

var (
	// Verbose is a global mode, when set at true, it will override the
	// configuration of the connections.
	Verbose bool

	// Debug is a global mode, when set at true, it will override the
	// configuration of the connections.
	Debug bool

	// m is the mutex that manage the pool of sqlx wrapper.
	m sync.Mutex

	// pool of sqlx wrapper.
	cp = make([]*sqlx.DB, 0, 5)

	// map of ...
	cm = make(map[string]uint, 5)
)

type (
	// MapMapper is the prototype to map a map result.
	MapMapper func(map[string]interface{})

	// SliceMapper is the prototype to map a slice result.
	SliceMapper func([]interface{})

	// Connection is a connection to an database.
	Connection interface {
		DB() *sqlx.DB
		Copy() Connection
		Connect() error
		LastError() error
		Ping() error
		MustPing()
		Close() error
		SetLogger(out *log.Logger, err *log.Logger)

		// Statements
		Exec(Stmt) (sql.Result, error)

		// Queries
		SelectMap(Stmt, MapMapper) (int, error)
		SelectSlice(Stmt, SliceMapper) (int, error)
		SelectMapRow(Stmt, MapMapper) (int, error)
		SelectSliceRow(Stmt, SliceMapper) (int, error)

		QueryMap(string, MapMapper, ...interface{}) (int, error)
		QuerySlice(string, SliceMapper, ...interface{}) (int, error)
		QueryMapRow(string, MapMapper, ...interface{}) (int, error)
		QuerySliceRow(string, SliceMapper, ...interface{}) (int, error)

		// Context
		Context(...string) Connection
		Done()
		HasContext() bool
		RunContext(...ContextFunc) error

		// Tx
		IsTx() bool
		Tx(...sql.IsolationLevel) (Connection, error)
		Commit() error
		Rollback() error
		RunTx(...TxFunc) error
	}
)

// Close closes active connections in the pool.
// it closes any profiler running in background.
func Close() {
	m.Lock()
	defer m.Unlock()
	for _, dbx := range cp {
		dbx.Close()
	}
}

// Open opens a database from default environement.
func Open(logger ...*log.Logger) (Connection, error) {
	return open(nil, false, "", nil, logger)
}

// OpenOnce opens a new connection or return the existing one.
func OpenOnce(logger ...*log.Logger) (Connection, error) {
	return open(nil, true, "", nil, logger)
}

// OpenWith opens a database with given connection.
func OpenWith(dbx *sqlx.DB, logger ...*log.Logger) (Connection, error) {
	return openWith(dbx, logger)
}

// OpenEnv opens a database from given environment.
func OpenEnv(env string, logger ...*log.Logger) (Connection, error) {
	return open(nil, false, env, nil, logger)
}

// OpenOnceEnv opens a database from a given environment or return an existing one.
func OpenOnceEnv(env string, logger ...*log.Logger) (Connection, error) {
	return open(nil, true, env, nil, logger)
}

// OpenEnviron opens a database from a given environ.
func OpenEnviron(e Environment, logger ...*log.Logger) (Connection, error) {
	return openEnviron(e, false, logger)
}

// OpenOnceEnviron opens a database from a given environ or return existing one.
func OpenOnceEnviron(e Environment, logger ...*log.Logger) (Connection, error) {
	return openEnviron(e, true, logger)
}

// -------------------------------------------------

// openEnv open a new connection through environment variables.
func openEnv(env string, once bool, logger []*log.Logger) (*db, error) {
	return open(nil, once, env, nil, logger)
}

// openEnviron open a new connection with a given Environ.
func openEnviron(e Environment, once bool, logger []*log.Logger) (*db, error) {
	return open(&e, once, "", nil, logger)
}

// openWith open a new connection with a given sqlx.DB connection.
func openWith(dbx *sqlx.DB, logger []*log.Logger) (*db, error) {
	return open(nil, false, "", dbx, logger)
}

// open a new connection based on input.
func open(e *Environment, once bool, env string, dbx *sqlx.DB, logger []*log.Logger) (*db, error) {
	var cfg Environment

	if e == nil {
		cfg.Boot()
		if env != "" {
			cfg.Load(env)
		}
	} else {
		cfg = *e
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	conn := &db{
		env: cfg,
		m:   &sync.Mutex{},
	}

	if once {
		m.Lock()

		m.Unlock()
	}
	conn.id, conn.dbx = createNewConnection(once, cfg.Alias)

	if conn.hasVerbose() {
		if n := len(logger); n > 0 {
			conn.logOut = logger[0]

			if n > 1 {
				conn.logErr = logger[1]
			}
		} else {
			prefix := fmt.Sprintf(
				"(\033[95;1m%s\033[0m:%s@%s:%s - \033[4m%s\033[0m)",
				conn.env.Driver,
				conn.env.User,
				conn.env.Host,
				conn.env.Port,
				conn.env.Alias,
			)

			conn.logOut = log.New(os.Stdout, fmt.Sprintf("%s ➜  ", prefix), 0)
			conn.logErr = log.New(os.Stderr, fmt.Sprintf("%s \033[91m➜ \033[1mERROR: \033[0m ", prefix), 0)
		}

		if conn.env.DSN == "" {
			conn.logOut.Print("configured and ready")
			conn.logOut.Printf("setting: %d MaxIdle | %d MaxOpen | %s MaxLifetime", conn.env.MaxIdle, conn.env.MaxOpen, conn.env.MaxLifetime.String())
		}
	}

	if conn.env.Autoconnect {
		conn.Connect()
	}

	return conn, nil
}

// createNewConnection in the pool "cp".
func createNewConnection(once bool, alias string, dbx ...*sqlx.DB) (uint, **sqlx.DB) {
	m.Lock()
	defer m.Unlock()

	if once {
		if i, ok := cm[alias]; ok && cp[i] != nil {
			return i, &cp[i]
		}
	}

	cid := uint(len(cp))
	if len(dbx) > 0 {
		cp = append(cp, dbx[0])
	} else {
		cp = append(cp, nil)
	}
	if once {
		cm[alias] = cid
	}

	return cid, &cp[cid]
}
