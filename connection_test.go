// Copyright © 2019 Alexandre Kovac <contact@kovacou.fr>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package database

import (
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

// TestMysql determine if we run the integration tests.
var TestMysql = false

func init() {
	_, TestMysql = os.LookupEnv("DATABASE_DSN")

	if TestMysql {
		setupEnvironment()
	}
}

// create & drop schema queries w/ fixtures insert queries.
const (
	createSchema = `
		CREATE TABLE articles (
			id integer,
			author_id integer,
			title text,
			content text,
			created_at datetime default NOW()
		);

		CREATE TABLE authors (
			id integer,
			firstname text null,
			lastname text null,
			nickname text
		);
	`

	dropSchema = `
		DROP TABLE IF EXISTS articles;
		DROP TABLE IF EXISTS authors;
	`

	createFixtures = `
		INSERT INTO authors(id, firstname, lastname, nickname) VALUES(1, "Alexandre", "Smith", "Kovacou");
		INSERT INTO authors(id, firstname, lastname, nickname) VALUES(2, "Sebastian", NULL, "Saucisse");

		INSERT INTO articles(id, author_id, title, content) VALUES(1, 1, "go-database package introduced", "See github.com/kovacou/go-database");
		INSERT INTO articles(id, author_id, title, content) VALUES(2, 2, "TEST", "This is a beautiful test, sorry.");
	`
)

// Article test struct.
type Article struct {
	ID        uint64
	AuthorID  uint64
	Title     string
	Content   string
	CreatedAt time.Time
}

// Author test struct.
type Author struct {
	ID        uint64
	Nickname  string
	Firstname *string
	Lastname  *string
}

// setupEnvironment will setup the test environment.
func setupEnvironment() {
	conn := sqlx.MustOpen("mysql", os.Getenv("DATABASE_DSN"))
	conn.Exec(dropSchema)
	conn.Exec(createSchema)
	conn.Exec(createFixtures)
}

func TestConnect(t *testing.T) {
	if !TestMysql {
		return
	}

	// Working case
	{
		conn, err := OpenEnv("DSN")
		assert.NoError(t, err)
		assert.NotNil(t, conn)
		assert.Nil(t, conn.Connect())
		assert.Nil(t, conn.Close())
	}

	// Not working case
	{
		conn, err := OpenEnv("NOK")
		assert.NoError(t, err)
		assert.NotNil(t, conn)
		assert.Error(t, conn.Connect())
	}
}

func TestPing(t *testing.T) {
	if !TestMysql {
		return
	}

	// Working case
	{
		conn, err := OpenEnv("DSN")
		assert.NoError(t, err)
		assert.NotNil(t, conn)
		assert.Nil(t, conn.Ping())
		assert.Nil(t, conn.Close())
	}

	// Not working case
	{
		conn, err := OpenEnv("NOK")
		assert.NoError(t, err)
		assert.NotNil(t, conn)
		assert.NotNil(t, conn.Ping())
	}
}

func TestMustPing(t *testing.T) {
	if !TestMysql {
		return
	}

	// Working case
	{
		conn, err := OpenEnv("DSN")
		assert.NoError(t, err)
		assert.NotNil(t, conn)
		assert.NotPanics(t, conn.MustPing)
		assert.Nil(t, conn.Close())
	}

	// Not working case
	{
		conn, err := OpenEnv("NOK")
		assert.NotNil(t, conn)
		assert.NoError(t, err)
		assert.Panics(t, conn.MustPing)
	}
}
