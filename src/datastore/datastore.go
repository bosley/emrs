package datastore

import (
	"time"
)

type InterfacePanel struct {
	ServerDb ServerStore
	UserDb   UserStore
	AssetDb  AssetStore
	Handler  ControlHandle
}

type ControlHandle interface {
	Close()
}

type ServerStore interface {
	LoadIdentity() *string
	StoreIdentity(identity string) error
}

type UserStore interface {
	GetAuthHash(username string) *string
	AddUser(name string, password string) error
	UpdatePassword(username string, password string) error
	DeleteUser(username string) bool
}

type Asset struct {
	Id          int
	Name        string
	Description string
}

type AssetStore interface {
	GetAsset(name string) *Asset
	GetAssets() []Asset
	AddAsset(name string, description string) error
	DeleteAsset(name string) error
	UpdateAsset(name string, description string) error
}

type Event struct {
	Id          int
	Received    time.Time
	OriginAsset int
	Data        string
}

func New(path string) (InterfacePanel, error) {
	c, err := newController(path)
	if err != nil {
		return InterfacePanel{}, err
	}
	return InterfacePanel{
		ServerDb: c,
		UserDb:   c,
		AssetDb:  c,
		Handler:  c,
	}, nil
}
