package core

import (
  "time"
  "sync"
  "sync/atomic"
)

type Service interface {
  Start() error
  Stop() error
  ShutdownAlert(time.Duration)
  Alive() bool
}

type Core struct {
  running   atomic.Bool
  stats      *stats
  serviceMgr *serviceManager
  wg         *sync.WaitGroup
}

type stats struct {
  start   time.Time
}

func New() *Core {
  return &Core{
    stats: nil,
    serviceMgr: newServiceManager(),
    wg: new(sync.WaitGroup),
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

func (e *Core) AddService(name string, service Service) error {
  if e.running.Load() {
    return ErrNotPermittedOnline
  }

  if err := e.serviceMgr.add(name, service); err != nil {
    return err
  }

  return nil
}
