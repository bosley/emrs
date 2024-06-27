package main

import (
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

var buildInfo = core.BuildInfo {
  Major: 0,
  Minor: 0,
  Patch: 0,
  Release: false,
}

func main() {

	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
				})))

	releaseMode := flag.Bool("release", buildInfo.Release, "Turn on debug mode")
	selectedConfig := flag.String("config", configName, "Use specified config file")

	flag.Parse()

  buildInfo.Release = *releaseMode

	cfg, err := LoadConfig(*selectedConfig)

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

	appCore := core.New(buildInfo, dbip)

	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	if err := appCore.AddService(
    "webui", 
    webui.New(
      appCore,
      cfg.GetAddress(),
      cfg.Assets,
      cert)); err != nil {
		slog.Error("Failed to add webui service to application core", "error", err.Error())
		panic("failed to create web ui service")
  }

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

