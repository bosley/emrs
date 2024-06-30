package datastore

import (
	"errors"
	"log/slog"
)

var ErrorSignalMapExists = errors.New("signalmap name already exists")

func (c *controller) GetSignalMap(id int) *SignalMap {
	return c.retrieveSignalMap(id)
}

func (c *controller) GetSignalMaps() []SignalMap {

	slog.Debug("retreiving all signalmaps")
	result := make([]SignalMap, 0)

	rows, err := c.db.Query(signalmaps_select_all)
	if err != nil {
		slog.Error("failed to select signalmaps", "err", err.Error())
		return result
	}
	defer rows.Close()

	for rows.Next() {
		entry := SignalMap{0, 0, 0}
		err = rows.Scan(&entry.Id, &entry.SignalId, &entry.ActionId)
		if err != nil {
			slog.Error("error scanning signalmap entry", "error", err.Error())
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

func (c *controller) AddSignalMap(signalId int, actionId int) error {

	slog.Debug("checking to see if signalmap exists")

	if u := c.findPair(signalId, actionId); u != nil {
		return ErrorSignalMapExists
	}

	tx, err := c.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(signalmaps_create)
	if err != nil {
		slog.Error("Error preparing signalmap-create", "err", err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(signalId, actionId)
	err = tx.Commit()
	if err != nil {
		slog.Error("Error tx commit", "err", err.Error())
		return err
	}
	return nil
}

func (c *controller) DeleteSignalMap(id int) error {

	slog.Debug("deleting signalmap", "id", id)

	tx, err := c.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(signalmaps_delete)
	if err != nil {
		slog.Error("Error preparing signalmaps-delete", "err", err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	err = tx.Commit()
	if err != nil {
		slog.Error("Error tx commit", "err", err.Error())
		return err
	}
	return nil
}

func (c *controller) UpdateSignalMap(id int, signalId int, actionId int) error {

	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(signalmaps_update)
	if err != nil {
		slog.Error("Error preparing signalmap-update", "err", err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(signalId, actionId, id)
	err = tx.Commit()
	if err != nil {
		slog.Error("Error tx commit", "err", err.Error())
		return err
	}
	return nil
}

const signalmaps_create = `insert into signalmaps (id, signalId, actionId) values (NULL, ?, ?)`
const signalmaps_get = `select id, signalId, actionId from signalmaps where id = ?`
const signalmaps_getpair = `select id, signalId, actionId from signalmaps where signalId = ? and actionId = ?`
const signalmaps_delete = `delete from signalmaps where id = ?`
const signalmaps_select_all = `select id, signalId, actionId from signalmaps where id >= 0`
const signalmaps_update = `update signalmaps set signalId = ?, actionId = ? where id = ?`

func (c *controller) retrieveSignalMap(id int) *SignalMap {

	stmt, err := c.db.Prepare(signalmaps_get)

	if err != nil {
		slog.Error("Error retreiving signalmap", "signalmap", id, "err", err.Error())
		return nil
	}
	defer stmt.Close()

	a := SignalMap{}

	err = stmt.QueryRow(id).Scan(&a.Id, &a.SignalId, &a.ActionId)
	if err == nil {
		return &a
	}
	slog.Error("Error retreiving signalmap post-query", "signalmap", id, "err", err.Error())
	return nil
}

func (c *controller) findPair(sig int, act int) *SignalMap {

	stmt, err := c.db.Prepare(signalmaps_getpair)

	if err != nil {
		slog.Error("Error retreiving signalmap", "s", sig, "a", act, "err", err.Error())
		return nil
	}
	defer stmt.Close()

	a := SignalMap{}

	err = stmt.QueryRow(sig, act).Scan(&a.Id, &a.SignalId, &a.ActionId)
	if err == nil {
		return &a
	}
	slog.Error("Error retreiving signalmap", "s", sig, "a", act, "err", err.Error())
	return nil
}
