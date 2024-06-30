package datastore

import (
	"errors"
	"log/slog"
)

var ErrorSignalExists = errors.New("signal name already exists")

func (c *controller) GetSignal(name string) *Signal {
	return c.retrieveSignal(name)
}

func (c *controller) GetSignals() []Signal {

	slog.Debug("retreiving all signals")
	result := make([]Signal, 0)

	rows, err := c.db.Query(signals_select_all)
	if err != nil {
		slog.Error("failed to select signals", "err", err.Error())
		return result
	}
	defer rows.Close()

	for rows.Next() {
		entry := Signal{0, "", "", ""}
		err = rows.Scan(&entry.Id, &entry.Name, &entry.Description, &entry.Triggers)
		if err != nil {
			slog.Error("error scanning signal entry", "error", err.Error())
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

func (c *controller) AddSignal(name string, description string, triggers string) error {

	slog.Debug("checking to see if signal exists")

	if u := c.retrieveSignal(name); u != nil {
		return ErrorSignalExists
	}

	tx, err := c.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(signals_create)
	if err != nil {
		slog.Error("Error preparing signal-create", "err", err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, description, triggers)
	err = tx.Commit()
	if err != nil {
		slog.Error("Error tx commit", "err", err.Error())
		return err
	}
	return nil
}

func (c *controller) DeleteSignal(name string) error {

	slog.Debug("deleting signal", "name", name)

	tx, err := c.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(signals_delete)
	if err != nil {
		slog.Error("Error preparing signals-delete", "err", err.Error())
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

func (c *controller) UpdateSignal(original string, name string, desc string, triggers string) error {

	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(signals_update)
	if err != nil {
		slog.Error("Error preparing signal-update", "err", err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, desc, triggers, original)
	err = tx.Commit()
	if err != nil {
		slog.Error("Error tx commit", "err", err.Error())
		return err
	}
	return nil
}

const signals_create = `insert into signals (id, name, description, triggers) values (NULL, ?, ?, ?)`
const signals_get = `select id, name, description, triggers from signals where name = ?`
const signals_delete = `delete from signals where name = ?`
const signals_select_all = `select id, name, description, triggers from signals where id >= 0`
const signals_update = `update signals set name = ?, description = ?, triggers = ? where name = ?`

func (c *controller) retrieveSignal(name string) *Signal {

	stmt, err := c.db.Prepare(signals_get)

	if err != nil {
		slog.Error("Error retreiving signal", "signal", name, "err", err.Error())
		return nil
	}
	defer stmt.Close()

	a := Signal{}

	err = stmt.QueryRow(name).Scan(&a.Id, &a.Name, &a.Description)
	if err == nil {
		return &a
	}
	slog.Error("Error retreiving signal post-query", "signal", name, "err", err.Error())
	return nil
}
