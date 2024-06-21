package webui

import (
  "time"
  "log/slog"
  "emrs/core"
	"crypto/tls"
  "sync/atomic"
)

func New(
  appCore *core.Core,
  address string,
  cert tls.Certificate) *controller {
  return &controller {
    appCore: appCore,
    address: address,
    tlsCert: cert,
  }
}

type controller struct {
  appCore *core.Core
  address string
  tlsCert tls.Certificate
  running atomic.Bool
}

func (c *controller) Start() error {
  if c.running.Load() {
    return nil
  }
  c.running.Store(true)
  slog.Info("webui started")

  slog.Warn("TODO: START THE WEB SERVER")

  return nil
}

func (c *controller) Stop() error {
  return nil
}

func (c *controller) ShutdownAlert(time.Duration) {
  slog.Warn("webui received shutdown alert")
}

func (c *controller) Alive() bool {
  return c.running.Load()
}
