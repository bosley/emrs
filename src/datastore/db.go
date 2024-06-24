package datastore

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log/slog"
	"sync/atomic"
)

const (
	maxCons = 255
)

type controller struct {
	running atomic.Bool
	newDb   bool
	db      *sql.DB
}

func (c *controller) IsNew() bool {
	return c.newDb
}

func (c *controller) Close() {
	if c.db != nil {
		c.db.Close()
	}
}

func (c *controller) UpdateIdentity(identity string) bool {

	return false
}

const db_table_create_identity = `create table identity (
  id integer not null primary key,
  data text
)`

const db_table_create_users = `create table users (
  id integer not null primary key,
  username text,
  authhash text,
  UNIQUE(username)
)`

const db_table_create_assets = `create table assets (
  id integer not null primary key,
  owner int,
  name text,
  type int,
  description text,
  UNIQUE(name),
  FOREIGN KEY (owner) REFERENCES users(id)
)`

const db_contains_table = `select name from sqlite_master where type = 'table' and name = ?`

func newController(path string) (*controller, error) {
	slog.Debug("db_open")

	var c controller

	c.running.Store(false)

	const options = "?_journal_mode=WAL"
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s%s", path, options))
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxCons)

	c.db = db

	type tcs struct {
		name string
		stmt string
	}

	for _, table := range []tcs{
		tcs{"identity", db_table_create_identity},
		tcs{"users", db_table_create_users},
		tcs{"assets", db_table_create_assets},
	} {

		if err := db_ensure_table_exists(c.db, table.name, table.stmt); err != nil {
			slog.Error("error setting up table", "name", table.name)
			c.db.Close()
			return nil, err
		}
	}
	return &c, nil
}

func db_does_table_exist(db *sql.DB, table string) (bool, error) {

	slog.Debug("db_does_table_exist")

	stmt, err := db.Prepare(db_contains_table)

	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var name string
	err = stmt.QueryRow(table).Scan(&name)
	if err == nil {
		return true, nil
	}
	return false, nil
}

func db_ensure_table_exists(db *sql.DB, table string, creation_stmt string) error {

	slog.Debug("db_ensure_table_exists")

	exists, err := db_does_table_exist(db, table)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	_, err = db.Exec(creation_stmt)
	if err != nil {
		return err
	}

	return nil
}
