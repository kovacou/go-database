// Copyright © 2019 Alexandre Kovac <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/kovacou/go-database/builder"
)

// Stmt is the representation of an statement or query (SELECT, UPDATE, & DELETE)
type Stmt interface {
	String() string
	Args() []any
}

// SelectMap run an SELECT query to fetch multiple results using a map mapper.
func (conn *db) SelectMap(stmt Stmt, mapper MapMapper) (rowsReturned int, err error) {
	return conn.runMap(stmt, mapper)
}

// SelectMapRow run an SELECT query to fetch a single result using a map mapper.
func (conn *db) SelectMapRow(stmt Stmt, mapper MapMapper) (rowsReturned int, err error) {
	return conn.runMapRow(stmt, mapper)
}

// SelectSlice run an SELECT query to fetch multiple results using a slice mapper.
func (conn *db) SelectSlice(stmt Stmt, mapper SliceMapper) (rowsReturned int, err error) {
	return conn.runSlice(stmt, mapper)
}

// SelectSliceRow run an SELECT query to fetch a single result using a slice mapper.
func (conn *db) SelectSliceRow(stmt Stmt, mapper SliceMapper) (rowsReturned int, err error) {
	return conn.runSliceRow(stmt, mapper)
}

// Exec run a statement.
func (conn *db) Exec(stmt Stmt) (res sql.Result, err error) {
	if err := conn.Connect(); err != nil {
		return nil, err
	}

	var t time.Time
	if conn.hasProfiling() {
		t = time.Now()
	}

	if conn.tx != nil {
		res, err = conn.tx.Exec(stmt.String(), stmt.Args()...)
	} else {
		res, err = (*conn.dbx).Exec(stmt.String(), stmt.Args()...)
	}

	conn.profilingStmt(stmt, err, t)
	return
}

// QuerySlice run an SELECT query to fetch a multiple results using a slice mapper.
func (conn *db) QuerySlice(query string, mapper SliceMapper, args ...any) (rowsReturned int, err error) {
	return conn.runSlice(builder.NewQuery(query, args...), mapper)
}

// QuerySliceRow run an SELECT query to fetch a single result using a slice mapper.
func (conn *db) QuerySliceRow(query string, mapper SliceMapper, args ...any) (rowsReturned int, err error) {
	return conn.runSliceRow(builder.NewQuery(query, args...), mapper)
}

// QueryMap run an SELECT query to fetch multiple results using a map mapper.
func (conn *db) QueryMap(query string, mapper MapMapper, args ...any) (rowsReturned int, err error) {
	return conn.runMap(builder.NewQuery(query, args...), mapper)
}

// QueryMapRow run an SELECT query to fetch a single result using a map mapper.
func (conn *db) QueryMapRow(query string, mapper MapMapper, args ...any) (rowsReturned int, err error) {
	return conn.runMapRow(builder.NewQuery(query, args...), mapper)
}

// runMap run stmt with a multiple results expected and mapped with a MapMapper.
func (conn *db) runMap(stmt Stmt, mapper MapMapper) (rowsReturned int, err error) {
	if err = conn.Connect(); err != nil {
		return
	}

	var (
		stmtx *sqlx.Stmt
		rows  *sqlx.Rows
		t     time.Time
	)

	if conn.hasProfiling() {
		t = time.Now()
	}

	stmtx, err = preparex(conn, stmt)
	if err == nil {
		defer stmtx.Close()
		rows, err = stmtx.Queryx(stmt.Args()...)
		if err == nil {
			defer rows.Close()

			row := map[string]any{}
			for rows.Next() {
				err = rows.MapScan(row)
				if err != nil {
					break
				}
				mapper(row)
				rowsReturned++
			}
		} else if errors.Is(err, sql.ErrNoRows) {
			if !conn.env.ErrorNoRows {
				err = nil
			}
		}
	}

	if err != nil && conn.hasVerbose() {
		conn.logErr.Println(err.Error())
	}

	conn.profilingStmt(stmt, err, t)
	return
}

// runMapRow run stmt with a single result expected and mapped with a MapMapper.
func (conn *db) runMapRow(stmt Stmt, mapper MapMapper) (rowsReturned int, err error) {
	if err = conn.Connect(); err != nil {
		return
	}

	var (
		stmtx  *sqlx.Stmt
		t      time.Time
		values = map[string]any{}
	)

	if conn.hasProfiling() {
		t = time.Now()
	}

	stmtx, err = preparex(conn, stmt)
	if err == nil {
		defer stmtx.Close()

		err = stmtx.QueryRowx(stmt.Args()...).MapScan(values)
		if err == nil {
			mapper(values)
			rowsReturned = 1
		} else if errors.Is(err, sql.ErrNoRows) {
			if !conn.env.ErrorNoRows {
				err = nil
			}
		}
	}

	if err != nil && conn.hasVerbose() {
		conn.logErr.Println(err.Error())
	}

	conn.profilingStmt(stmt, err, t)
	return
}

// runSlice run stmt with a multiple results and mapped with a SliceMapper.
func (conn *db) runSlice(stmt Stmt, mapper SliceMapper) (rowsReturned int, err error) {
	if err = conn.Connect(); err != nil {
		return
	}

	var (
		stmtx  *sqlx.Stmt
		rows   *sqlx.Rows
		values []any
		t      time.Time
	)

	if conn.hasProfiling() {
		t = time.Now()
	}

	stmtx, err = preparex(conn, stmt)
	if err == nil {
		defer stmtx.Close()
		rows, err = stmtx.Queryx(stmt.Args()...)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				values, err = rows.SliceScan()
				if err != nil {
					break
				}

				mapper(values)
				rowsReturned++
			}
		} else if errors.Is(err, sql.ErrNoRows) {
			if !conn.env.ErrorNoRows {
				err = nil
			}
		}
	}

	if err != nil && conn.hasVerbose() {
		conn.logErr.Println(err.Error())
	}

	conn.profilingStmt(stmt, err, t)
	return
}

// runSliceRow run stmt with a single result expected and mapped with a SliceMapper.
func (conn *db) runSliceRow(stmt Stmt, mapper SliceMapper) (rowsReturned int, err error) {
	if err = conn.Connect(); err != nil {
		return
	}

	var (
		stmtx  *sqlx.Stmt
		values []any
		t      time.Time
	)

	if conn.hasProfiling() {
		t = time.Now()
	}

	stmtx, err = preparex(conn, stmt)
	if err == nil {
		defer stmtx.Close()
		if values, err = stmtx.QueryRowx(stmt.Args()...).SliceScan(); err == nil {
			mapper(values)
			rowsReturned = 1
		} else if errors.Is(err, sql.ErrNoRows) {
			if !conn.env.ErrorNoRows {
				err = nil
			}
		}
	}

	if err != nil && conn.hasVerbose() {
		conn.logErr.Println(err.Error())
	}

	conn.profilingStmt(stmt, err, t)
	return
}

// preparex will prepare a query based on the given connection.
func preparex(conn *db, stmt Stmt) (*sqlx.Stmt, error) {
	if conn.tx != nil {
		return conn.tx.Preparex(stmt.String())
	}

	return (*conn.dbx).Preparex(stmt.String())
}

// profilingStmt store into the context the Stmt and store
func (conn *db) profilingStmt(stmt Stmt, err error, t time.Time) {
	if err != nil {
		fmt.Println("Erreur", err.Error())
		fmt.Println("Stmt", stmt.String())
	}

	if !conn.hasProfiling() {
		return
	}

	qs := &qs{
		end:     time.Now(),
		query:   stmt.String(),
		args:    stmt.Args(),
		ctxID:   conn.ctx.id,
		ctxFlag: conn.ctx.flag,
		start:   t,
	}

	conn.ctx.Push(qs)
	conn.profiler.Push(qs)
}
