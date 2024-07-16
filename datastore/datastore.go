package datastore

import (
	c "github.com/ostafen/clover/v2"
	d "github.com/ostafen/clover/v2/document"
	"log/slog"
	"os"
	"path/filepath"
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

type User struct {
	DisplayName string `clover:"displayname"`
	Hash        string `clover:"hash"`
	UiKey       string `clover:"uikey"`
	Ring        int    `clover:"ring"`
}

type Ds struct {
	users  *c.DB
	assets *c.DB
}

func Load(location string) (*Ds, error) {

	u, err := c.Open(filepath.Join(location, usersDb))
	if err != nil {
		return nil, err
	}

	a, err := c.Open(filepath.Join(location, assetDb))
	if err != nil {
		return nil, err
	}

	return &Ds{
		users:  u,
		assets: a,
	}, nil
}

func (ds *Ds) Close() {
	ds.users.Close()
	ds.assets.Close()
}

func SetupDisk(location string, user User) {

	slog.Info("setting up datastore", "dir", location)

	for _, x := range []string{
		usersDb,
		assetDb,
	} {

		newPath := filepath.Join(location, x)
		os.MkdirAll(newPath, 0755)
		err := createCollection(x, newPath)
		if err != nil {
			slog.Error("failed to create collections", "error", err.Error())
			panic("setup failure")
		}
	}

	slog.Info("forcing first user to be ring one", "user", user.DisplayName)
	user.Ring = RingOne

	createUser(user, filepath.Join(location, usersDb))
}

func createUser(user User, path string) {

	db, err := c.Open(path)
	if err != nil {
		slog.Error("failed to create user", "user", user.DisplayName, "error", err.Error())
		panic("setup failure")
	}

	defer db.Close()

	documentId, err := db.InsertOne(usersDb, user.ToDoc())
	if err != nil {
		slog.Error("failed to create user", user.DisplayName, "error", err.Error())
		panic("setup failure")
	}

	slog.Debug("new user created", "user", user.DisplayName, "doc-id", documentId)
}

func createCollection(name string, path string) error {

	slog.Info("create collection", "collection", name, "path", path)

	db, _ := c.Open(path)
	defer db.Close()

	if db == nil {
		panic("yep, it nil boss")
	}

	db.CreateCollection(name)

	return nil
}

func (u User) ToDoc() *d.Document {
	doc := d.NewDocument()
	doc.Set("displayname", u.DisplayName)
	doc.Set("hash", u.Hash)
	doc.Set("uikey", u.UiKey)
	doc.Set("ring", u.Ring)
	return doc
}
