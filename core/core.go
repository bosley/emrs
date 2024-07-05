package core

import (
	"emrs/badger"
	"errors"
	"log/slog"
	"maps"
	"sync"
	"sync/atomic"
	"time"
)

var ErrUpdating = errors.New("Core is busy updating")
var ErrNoMapping = errors.New("Asset does not have corresponding `onEvent` signal mapped")

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

type Event struct {
	Origin string `json:asset_origin`
	Data   string `json:data`
}

// Core EMRS object that maintains network configuration state,
// the pub/sub event bus and other internal workings
type Core struct {
	networkObservers []SnapshotReceiver

	badge   badger.Badge
	network *NetworkMap
	netmu   sync.Mutex

	updating atomic.Bool

	onEventMap map[string](chan Event)

	metrics Metrics
}

type Metrics struct {
	Created           time.Time     `json:time_created`
	SubmissionAttempt atomic.Uint64 `json:n_submission_attempt`
	SubmissionFailure atomic.Uint64 `json:n_submission_failure`
	SubmissionSuccess atomic.Uint64 `json:n_submission_success`
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

		badge:      badge,
		network:    nm,
		onEventMap: make(map[string](chan Event)),
	}

	core.updating.Store(false)

	core.metrics.Created = time.Now()
	return core, nil
}

// Submit an event from an asset. In the event that the asset.onEvent
// signal does not exist, then either the asset isn't valid, or the user
// doesn't have the onEvent signal for the given asset mapped to any actions
func (c *Core) SubmitEvent(originatingAsset string, data string) error {
	if c.updating.Load() {
		return ErrUpdating
	}

	signalNameOnEvent := formatAssetOnEventSignalName(originatingAsset)
	actionChannel, exists := c.onEventMap[signalNameOnEvent]
	if !exists {
		slog.Error("submission from valid un-mapped asset",
			"onEvent", signalNameOnEvent)
		return ErrNoMapping
	}
	actionChannel <- Event{
		Origin: originatingAsset,
		Data:   data,
	}
	return nil
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

	if c.updating.Load() {
		return ErrUpdating
	}

	c.updating.Store(true)
	defer c.updating.Store(false)

	slog.Info("updating network map")

	nm, err := NetworkMapFromTopo(topo)
	if err != nil {
		return err
	}

	c.netmu.Lock()
	defer c.netmu.Unlock()

	c.network = nm
	snapshot := c.getNetworkSnapshot()

	// onEvent signals are emitted by the core when an event for an
	// asset comes in so we need to setup the corresponding actions

	c.onEventMap = make(map[string](chan Event))

	for name, _ := range c.network.assets {
		onEvent := formatAssetOnEventSignalName(name)
		actions, ok := c.network.sigmap[onEvent]
		if !ok {
			continue
		}

		c.onEventMap[name] = make(chan Event)

		slog.Info("setting up onEvent for asset",
			"asset", name)

		for _, action := range actions {
			slog.Warn("NEET TO SETUP ACTION FOR EXEC",
				"action",
				action.Header.Name)

			//TODO: for file reading, we need to configure a runner
			//       to always read the file before interp

			//TODO:  for other type we just need to have the runner ready

			//TODO: hand the action  c.onEventMap[name] to consume on
		}
		return nil
	}

	// Inform all of the observers in parallel that we have
	// updated the network map
	iterate(
		c.networkObservers,
		func(iter Iter[SnapshotReceiver]) error {
			go iter.Value(&snapshot)
			return nil
		})

	slog.Info("network map updated")
	return nil
}
