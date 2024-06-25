package main

import (
	"crypto/tls"
	"emrs/core"
	"emrs/datastore"
	"emrs/webui"
	"flag"
	"log/slog"
	"os"
)

const (
	configName = "emrs_config.yaml"
)

func main() {

	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
				})))

	releaseMode := flag.Bool("release", false, "Turn on debug mode")

	flag.Parse()

	cfg, err := LoadConfig(configName)

	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
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

	dbip, err := datastore.New(cfg.Datastore)

	appCore := core.New(*releaseMode, dbip)

	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	setupServiceWebUi(
		appCore,
		cfg.GetAddress(),
		cert)

	if err := appCore.Start(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	appCore.Await()

	if err := appCore.Stop(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func setupServiceWebUi(appCore *core.Core, address string, cert tls.Certificate) {
	err := appCore.AddService("webui", webui.New(appCore, address, cert))
	if err != nil {
		slog.Error("Failed to add webui service to application core", "error", err.Error())
		panic("failed to create webui")
	}
}
