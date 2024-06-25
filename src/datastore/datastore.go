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
	IsNew() bool
	Close()
}

type ServerStore interface {
	LoadIdentity() *string
	StoreIdentity(identity string) error
}

type UserStore interface {
	Validate(username string, password string) bool

	AddUser(name string, password string) error

	UpdatePassword(username string, password string) error

	DeleteUser(username string) bool
}

const (
	AssetTypeRx = iota
	AssetTypeTx
	AssetTypeRxTx
)

type Asset struct {
	Id          int
	Name        string
	Type        int
	Description string
}

type AssetStore interface {
	GetAsset(Id string) *Asset

	GetAssets() []Asset

	AddAsset(name string, assetType int, description string) error

	DeleteAsset(Id string) error
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
