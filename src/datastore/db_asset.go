package datastore

import (
	"errors"
	"log/slog"
)

var ErrorAssetExists = errors.New("asset name already exists")

func (c *controller) GetAsset(name string) *Asset {
	return c.retrieveAsset(name)
}

func (c *controller) GetAssets() []Asset {

	slog.Debug("retreiving all assets")
	result := make([]Asset, 0)

	rows, err := c.db.Query(assets_select_all)
	if err != nil {
		slog.Error("failed to select assets", "err", err.Error())
		return result
	}
	defer rows.Close()

	for rows.Next() {
		entry := Asset{0, "", ""}
		err = rows.Scan(&entry.Id, &entry.Name, &entry.Description)
		if err != nil {
			slog.Error("error scanning asset entry", "error", err.Error())
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

func (c *controller) AddAsset(name string, description string) error {

	slog.Debug("checking to see if asset exists")

	if u := c.retrieveAsset(name); u != nil {
		return ErrorAssetExists
	}

	tx, err := c.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(assets_create)
	if err != nil {
		slog.Error("Error preparing asset-create", "err", err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, description)
	err = tx.Commit()
	if err != nil {
		slog.Error("Error tx commit", "err", err.Error())
		return err
	}
	return nil
}

func (c *controller) DeleteAsset(name string) error {

	slog.Debug("deleting asset", "name", name)

	tx, err := c.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(assets_delete)
	if err != nil {
		slog.Error("Error preparing assets-delete", "err", err.Error())
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

const assets_create = `insert into assets (id, name, description) values (NULL, ?, ?)`
const assets_get = `select id, name, description from assets where name = ?`
const assets_delete = `delete from assets where name = ?`
const assets_select_all = `select id, name, description from assets where id >= 0`

func (c *controller) retrieveAsset(name string) *Asset {

	stmt, err := c.db.Prepare(assets_get)

	if err != nil {
		slog.Error("Error retreiving user", "user", name, "err", err.Error())
		return nil
	}
	defer stmt.Close()

	a := Asset{}

	err = stmt.QueryRow(name).Scan(&a.Id, &a.Name, &a.Description)
	if err == nil {
		return &a
	}
	slog.Error("Error retreiving asset post-query", "asset", name, "err", err.Error())
	return nil
}
