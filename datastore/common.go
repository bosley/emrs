package datastore

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

const db_table_create_users = `create table users (
  id integer not null primary key,
  name text,
  hash text,
  key text,
  ring int,
  UNIQUE(name)
)`

const users_create = `insert into users (id, name, hash, key, ring) values (NULL, ?, ?, ?, ?)`
const users_get = `select name, hash, key, ring from users where name = ?`
const users_update = `update users set name = ?, hash = ?. key = ?, ring = ? where name = ?`
const users_delete = `delete from users where name = ?`
const users_load_owner = `select name, hash, key, ring from users where ring = ?`

const db_table_create_assets = `create table assets (
  id integer not null primary key,
  uuid text,
  name text
)`

const assets_create = `insert into assets (id, uuid, name) values (NULL, ?, ?)`
const assets_get = `select uuid, name from assets where uuid = ?`
const assets_update = `update assets set uuid = ?, name = ? where uuid = ?`
const assets_delete = `delete from assets where uuid = ?`
const assets_fetch = `select uuid, name from assets`

const db_contains_table = `select name from sqlite_master where type = 'table' and name = ?`

func db_does_table_exist(db *sql.DB, table string) (bool, error) {
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
