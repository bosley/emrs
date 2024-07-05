package core

import (
	"context"
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

type Event struct {
	Origin string `json:asset_origin`
	Data   string `json:data`
}

type EventReceiver func(Event)

type SnapshotReceiver func(*NetworkSnapshot)

// Core EMRS object that maintains network configuration state,
// the pub/sub event bus and other internal workings
type Core struct {
	networkObservers []SnapshotReceiver

	badge   badger.Badge
	network *NetworkMap
	netmu   sync.Mutex

	updating atomic.Bool

	actionChannels map[string](chan Event)

	ctx        context.Context
	procCancel context.CancelFunc
	procWg     *sync.WaitGroup

	metrics Metrics
}

type Metrics struct {
	Created            time.Time     `json:time_created`
	SubmissionAttempt  atomic.Uint64 `json:n_submission_attempt`
	SubmissionFailure  atomic.Uint64 `json:n_submission_failure`
	SubmissionSuccess  atomic.Uint64 `json:n_submission_success`
	CompletedProcesses int           `json:completed_processes`
	RunningProcesses   int           `json:running_processes`
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

		badge:          badge,
		network:        nm,
		actionChannels: make(map[string](chan Event)),
	}

	core.updating.Store(false)

	core.setupSync()

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
	actionChannel, exists := c.actionChannels[signalNameOnEvent]
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

func (c *Core) setupSync() {
	c.ctx = context.Background()
	c.ctx, c.procCancel = context.WithCancel(c.ctx)
	c.procWg = new(sync.WaitGroup)
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

	if c.metrics.RunningProcesses > 0 {
		slog.Debug("stopping previous action processors", "num", c.metrics.RunningProcesses)
		c.procCancel()
		c.procWg.Wait()
	}
	c.metrics.RunningProcesses = 0

	c.metrics.CompletedProcesses += c.metrics.RunningProcesses

	// Reset context for next run
	c.setupSync()

	c.network = nm
	snapshot := c.getNetworkSnapshot()

	c.actionChannels = make(map[string](chan Event))

	for signalName, actions := range c.network.sigmap {

		signal := c.network.signals[signalName]

		c.actionChannels[signalName] = make(chan Event)

		// Start a process to handle execution of every action
		// assigned to every signal
		// For a large network this might be a problem, but for a
		// self-hosted network of up-to a few hundred sensors/ reporters
		// this will be fine
		for _, action := range actions {
			proc := newProc(
				signal,
				action,
				c.ctx,
				c.actionChannels[signal.Header.Name])
			proc.exec(c.procWg)
			c.metrics.RunningProcesses += 1
		}
	}
	slog.Debug("action processes started",
		"count", c.metrics.RunningProcesses)

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
