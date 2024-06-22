package datastore

import (
	"sync/atomic"
)

type controller struct {
	running atomic.Bool
}

func (c *controller) UpdateIdentity(identity string) bool {

	return false
}
