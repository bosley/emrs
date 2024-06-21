package core

import (
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
	running    atomic.Bool
	stats      *stats
	serviceMgr *serviceManager
	wg         *sync.WaitGroup
	relMode    bool

	kt trigger
}

type stats struct {
	start time.Time
}

func New(releaseMode bool) *Core {
	return &Core{
		stats:      nil,
		serviceMgr: newServiceManager(),
		wg:         new(sync.WaitGroup),
		relMode:    releaseMode,
	}
}

func (e *Core) Start() error {
	if e.running.Load() {
		return ErrNotPermittedOnline
	}

	e.running.Store(true)

	if err := e.serviceMgr.start(); err != nil {
		e.running.Store(false)
		return err
	}

	e.kt = initReaperIntercept(
		e.wg,
		5*time.Second,
		func() {
			slog.Warn("Kill timer activated..")
			e.broadcastShutdownAlert()
		})

	return nil
}

func (e *Core) Await() {
	e.wg.Wait()
}

func (e *Core) Stop() error {
	if !e.running.Load() {
		return ErrNotPermittedOffline
	}

	if err := e.serviceMgr.stop(); err != nil {
		return err
	}

	return nil
}

func (e *Core) IsReleaseMode() bool {
	return e.relMode
}

func (e *Core) AddService(name string, service Service) error {
	if e.running.Load() {
		return ErrNotPermittedOnline
	}

	if err := e.serviceMgr.add(name, service); err != nil {
		return err
	}

	return nil
}

func (e *Core) broadcastShutdownAlert() {

	slog.Debug("TODO: Tell the services that we are about to shutdown")
}
