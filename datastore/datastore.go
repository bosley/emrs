package datastore

import (
	"log/slog"
	"os"
)

const (
	usersDb = "users"
	assetDb = "assets"
)

const (
	RingUnset = iota // Default int 0, thus declare as unset
	RingOne          // Ring one is the core, primary user (root) etc
	RingTwo          // Optional for later addition of non-owner users
)

type DataStore interface {
	AddAsset(asset Asset) bool
	GetAssets() []Asset
	UpdateAsset(asset Asset) bool
	RemoveAsset(id string) bool
	AssetExists(id string) bool

	GetOwner() (User, error)
	UpdateOwner(owner User) bool
	UpdateOwnerUiKey(key string) bool

	Close()
}

type User struct {
	DisplayName string
	Hash        string
	UiKey       string
	Ring        int
}

type Asset struct {
	Id          string
	DisplayName string
}

func Load(location string) (DataStore, error) {

	c, err := newController(location)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func SetupDisk(location string, user User) {

	slog.Info("setting up datastore", "dir", location)

	c, err := newController(location)
	if err != nil {
		slog.Error("Unable to open datastore", "location", location)
		os.Exit(1)
	}

	slog.Info("forcing first user to be ring one", "user", user.DisplayName)
	user.Ring = RingOne

	if err := c.createUser(user); err != nil {
		slog.Error("failed to setup initial user", "error", err.Error())
		os.Exit(1)
	}
}
