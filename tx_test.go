// Copyright © 2019 Alexandre Kovac <contact@kovacou.fr>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package database

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestTx(t *testing.T) {
	if !TestMysql {
		return
	}

	// Working case
	{
		conn, err := Open()
		assert.NoError(t, err)
		assert.NotNil(t, conn)
		tx, err := conn.Tx()
		assert.Nil(t, err)
		assert.NotNil(t, tx)
	}

	// Not working case
	{
		conn, _ := OpenEnv("NOK")

		tx, err := conn.Tx()
		assert.Error(t, err)
		assert.Nil(t, tx)
	}
}

func TestIsTx(t *testing.T) {
	{
		conn := &db{}
		assert.False(t, conn.IsTx())

		conn.tx = &sqlx.Tx{}
		assert.True(t, conn.IsTx())
	}

	if !TestMysql {
		return
	}

	// Working case
	{
		conn, err := Open()
		assert.NoError(t, err)
		assert.NotNil(t, conn)
		assert.False(t, conn.IsTx())
		tx, txErr := conn.Tx()
		assert.NoError(t, txErr)
		assert.True(t, tx.IsTx())
	}
}

func TestCommit(t *testing.T) {
	if !TestMysql {
		return
	}
}

func TestRollback(t *testing.T) {
	if !TestMysql {
		return
	}
}

func TestRunTx(t *testing.T) {
	if !TestMysql {
		return
	}
}
