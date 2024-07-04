package core

import (
	"emrs/badger"
	"maps"
	"sync"
)

// Base configuration information required to run a core
type Config struct {
	Identity string `json:identity`
	Network  Topo   `json:network`
}

// A snapshot copy of specific fields within the currently
// loaded network meant to be used for request validation
type NetworkSnapshot struct {

	// These are fully-qualified asset keys meaning
	// that they include the sector information
	Assets map[string]*Asset

	// All possible signals
	Signals map[string]*Signal

	// Signals that are mapped to actions
	// (unmapped actions are not important
	// for the runtime the application
	SignalMap map[string][]*Action
}

type SnapshotReceiver func(*NetworkSnapshot)

// Core EMRS object that maintains network configuration state,
// the pub/sub event bus and other internal workings
type Core struct {
	networkObservers []SnapshotReceiver

	badge   badger.Badge
	network *NetworkMap
	netmu   sync.Mutex
}

func New(config Config) (*Core, error) {

	nm, err := NetworkMapFromTopo(config.Network)
	if err != nil {
		return nil, err
	}

	badge, err := badger.DecodeIdentityString(config.Identity)
	if err != nil {
		return nil, err
	}

	core := &Core{
		networkObservers: make([]SnapshotReceiver, 0),

		badge:   badge,
		network: nm,
	}

	// TODO: core.setupEventBus

	return core, nil
}

// If the current network map is updated, then some
// functionality may require a new snapshot
func (c *Core) AddSnapshotReceiver(r SnapshotReceiver) *Core {
	c.networkObservers = append(c.networkObservers, r)
	return c
}

func (c *Core) getNetworkSnapshot() NetworkSnapshot {

	ns := NetworkSnapshot{
		Assets:    make(map[string]*Asset),
		Signals:   make(map[string]*Signal),
		SignalMap: make(map[string][]*Action),
	}

	maps.Copy(ns.Assets, c.network.assets)
	maps.Copy(ns.Signals, c.network.signals)
	maps.Copy(ns.SignalMap, c.network.sigmap)

	return ns
}

func (c *Core) UpdateNetworkMap(topo Topo) error {

	nm, err := NetworkMapFromTopo(topo)
	if err != nil {
		return err
	}

	c.netmu.Lock()
	defer c.netmu.Unlock()

	// TODO: Put event receiver pipeline from ingestion
	//        into a buffered mode so all submitted events
	//        are stored in memory temporarily

	// TODO: Stop the pub-sub functionality, remove all routes
	//        and all consumers. full reset

	c.network = nm
	snapshot := c.getNetworkSnapshot()

	// TODO: Restart/reconfigure pub-sub

	// TODO: Release the buffers onto the network. Keep storing events into buffer though
	//       until the moment the buffer is empty for the first time. Once its empty, disable
	//       the buffering. This will be to ensure that messages are delivered in the order
	//       that they are received by the endpoint

	// Inform all of the observers in parallel that we have
	// updated the network map
	iterate(
		c.networkObservers,
		func(iter Iter[SnapshotReceiver]) error {
			go iter.Value(&snapshot)
			return nil
		})

	return nil
}
