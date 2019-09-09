// Copyright © 2019 Alexandre Kovac <contact@kovacou.fr>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package database

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/kovacou/go-database/builder"
)

// Stmt is the representation of an statement or query (SELECT, UPDATE, & DELETE)
type Stmt interface {
	String() string
	//Bytes() []byte
	Args() []interface{}
}

// SelectMap run an SELECT query to fetch multiple results using a map mapper.
func (conn *db) SelectMap(stmt Stmt, mapper MapMapper) (rowsReturned int, err error) {
	return conn.runMap(stmt, mapper)
}

// SelectMap run an SELECT query to fetch a single result using a map mapper.
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
func (conn *db) Exec(stmt Stmt) (sql.Result, error) {
	if err := conn.Connect(); err != nil {
		return nil, err
	}

	if conn.tx != nil {
		return conn.tx.Exec(stmt.String(), stmt.Args()...)
	}

	return (*conn.dbx).Exec(stmt.String(), stmt.Args()...)
}

// QueryMap
func (conn *db) QueryMap(query string, mapper MapMapper, args ...interface{}) (rowsReturned int, err error) {
	return conn.runMap(builder.NewQuery(query, args), mapper)
}

// QueryMapRow
func (conn *db) QueryMapRow(query string, mapper MapMapper, args ...interface{}) (rowsReturned int, err error) {
	return conn.runMapRow(builder.NewQuery(query, args), mapper)
}

// runMap run stmt with a multiple results expected and mapped with a MapMapper.
func (conn *db) runMap(stmt Stmt, mapper MapMapper) (rowsReturned int, err error) {
	if err = conn.Connect(); err != nil {
		return
	}

	var (
		stmtx *sqlx.Stmt
		rows  *sqlx.Rows
	)

	stmtx, err = preparex(conn, stmt)

	if err == nil {
		defer stmtx.Close()
		rows, err = stmtx.Queryx(stmt.Args()...)
		if err == nil {
			defer rows.Close()

			row := map[string]interface{}{}
			for rows.Next() {
				err = rows.MapScan(row)
				if err != nil {
					break
				}
				mapper(row)
				rowsReturned++
			}
		}
	}
	if err != nil && conn.hasVerbose() {
		conn.log.Println(err.Error())
	}

	return
}

// runMapRow run stmt with a single result expected and mapped with a MapMapper.
func (conn *db) runMapRow(stmt Stmt, mapper MapMapper) (rowsReturned int, err error) {
	if err = conn.Connect(); err != nil {
		return
	}

	var (
		stmtx  *sqlx.Stmt
		values map[string]interface{}
	)

	stmtx, err = preparex(conn, stmt)
	if err == nil {
		defer stmtx.Close()

		err = stmtx.QueryRowx(stmt.Args()...).MapScan(values)
		if err == nil {
			mapper(values)
			rowsReturned = 1
		}
	}

	if err != nil && conn.hasVerbose() {
		conn.log.Println(err.Error())
	}

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
		values []interface{}
	)

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
		}
	}

	if err != nil && conn.hasVerbose() {
		conn.log.Println(err.Error())
	}

	return
}

// runSliceRow run stmt with a single result expected and mapped with a SliceMapper.
func (conn *db) runSliceRow(stmt Stmt, mapper SliceMapper) (rowsReturned int, err error) {
	if err = conn.Connect(); err != nil {
		return
	}

	var (
		stmtx  *sqlx.Stmt
		values []interface{}
	)

	stmtx, err = preparex(conn, stmt)
	if err == nil {
		defer stmtx.Close()
		if values, err = stmtx.QueryRowx(stmt.Args()...).SliceScan(); err == nil {
			mapper(values)
			rowsReturned = 1
		}
	}

	if err != nil && conn.hasVerbose() {
		conn.log.Println(err.Error())
	}

	return
}

// preparex will prepare a query based on the given connection.
func preparex(conn *db, stmt Stmt) (*sqlx.Stmt, error) {
	if conn.tx != nil {
		return conn.tx.Preparex(stmt.String())
	}

	return (*conn.dbx).Preparex(stmt.String())
}
