package core

import (
	"emrs/badger"
	ds "emrs/datastore"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

type Service interface {
	Start() error
	Stop() error
	ShutdownAlert(time.Duration)
	Alive() bool
}

type Core struct {
	badge      badger.Badge
	running    atomic.Bool
	stats      *stats
	serviceMgr *serviceManager
	wg         *sync.WaitGroup
	relMode    bool
	dbip       ds.InterfacePanel
	reqSetup   atomic.Bool

	kt trigger
}

type stats struct {
	start time.Time
}

func New(releaseMode bool, dbip ds.InterfacePanel) *Core {
	c := &Core{
		stats:      nil,
		serviceMgr: newServiceManager(),
		wg:         new(sync.WaitGroup),
		relMode:    releaseMode,
		dbip:       dbip,
	}
	c.setup()
	return c
}

func (c *Core) Start() error {
	if c.running.Load() {
		return ErrNotPermittedOnline
	}

	c.running.Store(true)

	if err := c.serviceMgr.start(); err != nil {
		c.running.Store(false)
		return err
	}

	c.kt = initReaperIntercept(
		c.wg,
		5*time.Second,
		func() {
			slog.Warn("Kill timer activated..")
			c.broadcastShutdownAlert()
		})

	return nil
}

func (c *Core) Await() {
	c.wg.Wait()
}

func (c *Core) Stop() error {
	if !c.running.Load() {
		return ErrNotPermittedOffline
	}

	if err := c.serviceMgr.stop(); err != nil {
		return err
	}

	return nil
}

func (c *Core) IsReleaseMode() bool {
	return c.relMode
}

func (c *Core) AddService(name string, service Service) error {
	if c.running.Load() {
		return ErrNotPermittedOnline
	}

	if err := c.serviceMgr.add(name, service); err != nil {
		return err
	}

	return nil
}

func (c *Core) broadcastShutdownAlert() {

	slog.Debug("TODO: Tell the services that we are about to shutdown")
}

func (c *Core) GetUserStore() ds.UserStore {
	return c.dbip.UserDb
}

func (c *Core) GetAssetStore() ds.AssetStore {
	return c.dbip.AssetDb
}

func (c *Core) RequiresSetup() bool {
	return c.reqSetup.Load()
}

func (c *Core) IndicateSetupComplete() {
	c.reqSetup.Store(false)
}

func (c *Core) GetSessionKey() []byte {
	return []byte(c.badge.Id())
}
