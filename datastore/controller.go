package datastore

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	_ "modernc.org/sqlite"
	"path/filepath"
	"sync/atomic"
)

const (
	maxCons = 255
	dbName  = "datastore.db"
)

var ErrorUserExists = errors.New("username already exists")

type controller struct {
	running atomic.Bool
	db      *sql.DB
}

func (c *controller) Close() {
	if c.db != nil {
		c.db.Close()
	}
}

func newController(path string) (*controller, error) {
	slog.Debug("db_open")

	path = filepath.Join(path, dbName)

	var c controller

	c.running.Store(false)

	const options = "?_journal_mode=WAL"
	db, err := sql.Open("sqlite", fmt.Sprintf("%s%s", path, options))
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

func (c *controller) createUser(user User) error {
	slog.Debug("checking to see if user exists")
	if u := c.retrieveUser(user.DisplayName); u != nil {
		return ErrorUserExists
	}
	slog.Debug("user does not yet exist")
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(users_create)
	if err != nil {
		slog.Error("error preparing user-create", "err", err.Error())
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		user.DisplayName,
		user.Hash,
		user.UiKey,
		user.Ring,
	)
	err = tx.Commit()
	if err != nil {
		slog.Error("error tx commit", "err", err.Error())
		return err
	}
	return nil
}

func (c *controller) AssetExists(id string) bool {

	slog.Debug(id)
	stmt, err := c.db.Prepare(assets_get)
	if err != nil {
		slog.Error("error retreiving asset", "id", id, "err", err.Error())
		return false
	}
	defer stmt.Close()
	u := Asset{}
	err = stmt.QueryRow(id).Scan(&u.Id, &u.DisplayName)
	if err == nil {
		return true
	}
	return false
}

func (c *controller) AddAsset(asset Asset) bool {
	if c.AssetExists(asset.DisplayName) {
		slog.Error("asset already exists")
		return false
	}
	tx, err := c.db.Begin()
	if err != nil {
		slog.Error("error creating tx", "error", err.Error())
		return false
	}
	stmt, err := tx.Prepare(assets_create)
	if err != nil {
		slog.Error("error preparing asset-create", "err", err.Error())
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		asset.Id,
		asset.DisplayName,
	)
	err = tx.Commit()
	if err != nil {
		slog.Error("error tx commit", "err", err.Error())
		return false
	}
	return true
}

func (c *controller) RemoveAsset(id string) bool {
	slog.Debug("deleting asset", "id", id)
	tx, err := c.db.Begin()
	if err != nil {
		slog.Error(err.Error())
		return false
	}
	stmt, err := tx.Prepare(assets_delete)
	if err != nil {
		slog.Error(err.Error())
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	err = tx.Commit()
	if err != nil {
		slog.Error(err.Error())
		return false
	}
	return true
}

func (c *controller) UpdateAsset(asset Asset) bool {
	slog.Debug("updatimg asset", "id", asset.Id)
	tx, err := c.db.Begin()
	if err != nil {
		slog.Error(err.Error())
		return false
	}
	stmt, err := tx.Prepare(assets_update)
	if err != nil {
		slog.Error(err.Error())
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(asset.Id, asset.DisplayName)
	err = tx.Commit()
	if err != nil {
		slog.Error(err.Error())
		return false
	}
	return true
}

func (c *controller) GetAssets() []Asset {
	slog.Debug("retrieve assets")
	result := make([]Asset, 0)
	rows, err := c.db.Query(assets_fetch)
	if err != nil {
		slog.Error(err.Error())
		return result
	}
	defer rows.Close()
	for rows.Next() {
		entry := Asset{"", ""}
		err = rows.Scan(&entry.Id, &entry.DisplayName)
		if err != nil {
			slog.Error(err.Error())
			return make([]Asset, 0)
		}
		result = append(result, entry)
	}
	err = rows.Err()
	if err != nil {
		slog.Error(err.Error())
		return make([]Asset, 0)
	}
	return result
}

func (c *controller) GetOwner() (User, error) {
	var u User
	stmt, err := c.db.Prepare(users_load_owner)
	if err != nil {
		return u, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(1).Scan(&u.DisplayName, &u.Hash, &u.UiKey, &u.Ring)
	if err == nil {
		return u, nil
	}
	return User{}, err
}

func (c *controller) UpdateOwnerUiKey(key string) bool {

	// Need to replace the key value for user where ring = 1
	panic("NOT YET IMPLEMENTED")
	return false
}

func (c *controller) UpdateOwner(owner User) bool {

	// Update user where ring = 1, make sure ring stays 1
	panic("NOT YET IMPLEMENTED")
	return true
}

// --

func (c *controller) retrieveUser(username string) *User {

	stmt, err := c.db.Prepare(users_get)

	if err != nil {
		slog.Error("error retreiving user", "user", username, "err", err.Error())
		return nil
	}
	defer stmt.Close()

	u := User{}

	err = stmt.QueryRow(username).Scan(&u.DisplayName, &u.Hash, &u.UiKey, &u.Ring)
	if err == nil {
		return &u
	}
	return nil
}
