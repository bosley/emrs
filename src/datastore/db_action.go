package datastore

import (
	"errors"
	"log/slog"
)

var ErrorActionExists = errors.New("action name already exists")

func (c *controller) GetAction(name string) *Action {
	return c.retrieveAction(name)
}

func (c *controller) GetActions() []Action {

	slog.Debug("retreiving all actions")
	result := make([]Action, 0)

	rows, err := c.db.Query(actions_select_all)
	if err != nil {
		slog.Error("failed to select actions", "err", err.Error())
		return result
	}
	defer rows.Close()

	for rows.Next() {
		entry := Action{0, "", "", ""}
		err = rows.Scan(&entry.Id, &entry.Name, &entry.Description, &entry.ExecutionInfo)
		if err != nil {
			slog.Error("error scanning action entry", "error", err.Error())
			return result
		}
		result = append(result, entry)
	}

	err = rows.Err()
	if err != nil {
		slog.Error("error post-query", "err", err.Error())
		return result
	}

	return result
}

func (c *controller) AddAction(name string, description string, execInfo string) error {

	slog.Debug("checking to see if action exists")

	if u := c.retrieveAction(name); u != nil {
		return ErrorActionExists
	}

	tx, err := c.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(actions_create)
	if err != nil {
		slog.Error("Error preparing action-create", "err", err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, description, execInfo)
	err = tx.Commit()
	if err != nil {
		slog.Error("Error tx commit", "err", err.Error())
		return err
	}
	return nil
}

func (c *controller) DeleteAction(name string) error {

	slog.Debug("deleting action", "name", name)

	tx, err := c.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(actions_delete)
	if err != nil {
		slog.Error("Error preparing actions-delete", "err", err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name)
	err = tx.Commit()
	if err != nil {
		slog.Error("Error tx commit", "err", err.Error())
		return err
	}
	return nil
}

func (c *controller) UpdateAction(original string, name string, desc string, execInfo string) error {

	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(actions_update)
	if err != nil {
		slog.Error("Error preparing action-update", "err", err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, desc, execInfo, original)
	err = tx.Commit()
	if err != nil {
		slog.Error("Error tx commit", "err", err.Error())
		return err
	}
	return nil
}

const actions_create = `insert into actions (id, name, description, executionInfo) values (NULL, ?, ?, ?)`
const actions_get = `select id, name, description, executionInfo from actions where name = ?`
const actions_delete = `delete from actions where name = ?`
const actions_select_all = `select id, name, description, executionInfo from actions where id >= 0`
const actions_update = `update actions set name = ?, description = ?, executionInfo = ? where name = ?`

func (c *controller) retrieveAction(name string) *Action {

	stmt, err := c.db.Prepare(actions_get)

	if err != nil {
		slog.Error("Error retreiving action", "action", name, "err", err.Error())
		return nil
	}
	defer stmt.Close()

	a := Action{}

	err = stmt.QueryRow(name).Scan(&a.Id, &a.Name, &a.Description)
	if err == nil {
		return &a
	}
	slog.Error("Error retreiving action post-query", "action", name, "err", err.Error())
	return nil
}
