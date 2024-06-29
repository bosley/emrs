package core

import (
	"emrs/badger"
	"log/slog"
)

func (c *Core) setup() {
	identity := c.dbip.ServerDb.LoadIdentity()
	if nil == identity {

		slog.Warn("No identity on file; creating a new identity and indicating setup required")

		// First time run, generate server identity and flag that we will need
		// to run any first time setup stuff
		badge, err := badger.New(badger.Config{
			Nickname: "EMRS Core",
		})
		if err != nil {
			slog.Error(err.Error())
			panic("badger error")
		}
		if err := c.dbip.ServerDb.StoreIdentity(badge.EncodeIdentityString()); err != nil {
			slog.Error(err.Error())
			panic("failed to store server identity")
		}
		c.badge = badge

		// Flag
		c.reqSetup.Store(true)

		slog.Debug("generated core identity - need to run setup", "id", badge.Id())

		return
	}

	// Not the first time, load the server identity
	badge, err := badger.DecodeIdentityString(*identity)
	if err != nil {
		slog.Error(err.Error())
		panic("failed to decode server identity")
	}
	c.badge = badge

	// Flag
	c.reqSetup.Store(false)

	slog.Debug("loaded core identity", "id", badge.Id())
}
