package datastore

import (
	"errors"
	"log/slog"
)

var ErrorUserExists = errors.New("username already exists")

type user struct {
	id   int
	name string
	auth string
}

func (c *controller) Validate(username string, password string) bool {

	u := c.retrieveUser(username)
	if u == nil {
		slog.Warn("unable to find user", "name", username)
		return false
	}

	slog.Warn("need to check pass")

	// TODO Use badger to check password here

	return password == u.auth
}

func (c *controller) AddUser(name string, password string) error {
	if u := c.retrieveUser(name); u != nil {
		return ErrorUserExists
	}

	tx, err := c.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(users_create)
	if err != nil {
		slog.Error("Error preparing user-create", "err", err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, password)
	err = tx.Commit()
	if err != nil {
		slog.Error("Error tx commit", "err", err.Error())
		return err
	}

	return nil
}

func (c *controller) UpdatePassword(username string, password string) error {
	return nil
}

func (c *controller) DeleteUser(username string) bool {
	return false
}

const users_create = `insert into users (id, username, authhash) values (NULL, ?, ?)`
const users_get = `select id, username, authhash from users where username = ?`
const users_update = `update users set username = ?, authhash = ? where username = ?`
const users_delete = `delete from users where username = ?`

func (c *controller) retrieveUser(username string) *user {

	stmt, err := c.db.Prepare(users_get)

	if err != nil {
		slog.Error("Error retreiving user", "user", username, "err", err.Error())
		return nil
	}
	defer stmt.Close()

	u := user{}

	err = stmt.QueryRow(username).Scan(&u.id, &u.name, &u.auth)
	if err == nil {
		return &u
	}
	slog.Error("Error retreiving user post-query", "user", username, "err", err.Error())
	return nil
}
