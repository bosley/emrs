package datastore

import (
	"time"
)

type InterfacePanel struct {
	ServerDb    ServerStore
	UserDb      UserStore
	AssetDb     AssetStore
	ActionDb    ActionStore
	SignalDb    SignalStore
	SignalMapDb SignalMapStore
	Handler     ControlHandle
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

type AssetStore interface {
	GetAsset(name string) *Asset
	GetAssets() []Asset
	AddAsset(name string, description string) error
	DeleteAsset(name string) error
	UpdateAsset(originalName string, name string, description string) error
}

type ActionStore interface {
	GetAction(name string) *Action
	GetActions() []Action
	AddAction(name string, description string, execInfo string) error
	DeleteAction(name string) error
	UpdateAction(originalName string, name string, description string, execInfo string) error
}

type SignalStore interface {
	GetSignal(name string) *Signal
	GetSignals() []Signal
	AddSignal(name string, description string, triggers string) error
	DeleteSignal(name string) error
	UpdateSignal(originalName string, name string, description string, triggers string) error
}

type SignalMapStore interface {
	GetSignalMap(id int) *SignalMap
	GetSignalMaps() []SignalMap
	AddSignalMap(signalId int, actionId int) error
	DeleteSignalMap(id int) error
	UpdateSignalMap(id int, signalId int, actionId int) error
}

type Asset struct {
	Id          int
	Name        string
	Description string
}

type Action struct {
	Id            int
	Name          string
	Description   string
	ExecutionInfo string
}

type Signal struct {
	Id          int
	Name        string
	Description string
	Triggers    string
}

type SignalMap struct {
	Id       int
	SignalId int
	ActionId int
}

type Event struct { // TODO: Add support for dumping all events to seperate database
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
		ServerDb:    c,
		UserDb:      c,
		AssetDb:     c,
		ActionDb:    c,
		SignalDb:    c,
		SignalMapDb: c,
		Handler:     c,
	}, nil
}
