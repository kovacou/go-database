// Copyright © 2019 Alexandre Kovac <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package database

import (
	"context"
	"database/sql"
)

// TxFunc handler.
type TxFunc func(Connection) error

// Tx copy the client and create a new transaction.
func (conn *db) Tx(level ...sql.IsolationLevel) (Connection, error) {
	var err error
	isolationLevel := sql.LevelDefault
	if len(level) > 0 {
		isolationLevel = level[0]
	}

	// try to connect to the database first.
	if err := conn.Connect(); err != nil {
		return nil, err
	}

	// create the transaction with the given isolation level.
	connTx := conn.copy()
	connTx.tx, err = (*conn.dbx).BeginTxx(context.Background(), &sql.TxOptions{
		Isolation: isolationLevel,
	})

	if err != nil {
		return nil, err
	}
	return connTx, nil
}

// RunTx run a bunch of TxFunc and handle the commit & rollback.
func (conn *db) RunTx(funcs ...TxFunc) (err error) {
	tx, err := conn.Tx()
	if err != nil {
		return
	}

	for _, f := range funcs {
		if err = f(tx); err != nil {
			tx.Rollback()
			return
		}
	}
	return tx.Commit()
}

// IsTx says if the current connection contains a Tx.
func (conn *db) IsTx() bool {
	return conn.tx != nil
}

// Commit the current transation.
func (conn *db) Commit() (err error) {
	if conn.IsTx() {
		err = conn.tx.Commit()
	}
	return
}

// Rollback the current transaction.
func (conn *db) Rollback() (err error) {
	if conn.IsTx() {
		err = conn.tx.Rollback()
	}
	return
}
