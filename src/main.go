package main

import (
  "os"
  "log/slog"
	"emrs/core"
	"crypto/tls"
)

func main() {

  slog.SetDefault(
		slog.New(
			slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
				})))

  cfg, err := LoadConfig("emrs.yaml")

  if err != nil {
    slog.Error(err.Error())
  }

  slog.Info("config",
	  "Home", cfg.Home,
	  "Hostname", cfg.Hostname,
	  "Port", cfg.Port,
	  "Key", cfg.Key,
	  "Cert", cfg.Cert,
	  "Datastore", cfg.Datastore)

  cert, err := cfg.LoadTlsCert()
  if err != nil {
    slog.Error("Failed to load TLS cert", "error", err.Error())
    os.Exit(1)
  }

  appCore := core.New()

  setupServiceWebUi(
    appCore,
    cfg.GetAddress(),
    cert)

  if err := appCore.Start(); err != nil {
    panic(err.Error())
  }

  appCore.Await()

  if err := appCore.Stop(); err != nil {
    panic(err.Error())
  }
}

func setupServiceWebUi(appCore* core.Core, address string, cert tls.Certificate) {

  slog.Debug("setup web ui", "address", address)


}
