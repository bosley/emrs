package main

import (
	"emrs/core"
	"errors"
	"flag"
	"log/slog"
	"os"
)

// --

const (
	defaultConfigName = "emrs.cfg"
)

func main() {

	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
				})))

	selectedConfig := flag.String("config", defaultConfigName, "Use specified config file")
	newConfig := flag.Bool("new", false, "Generate a template config file using [config] as the filename")
	overwriteConfig := flag.Bool("force", false, "Force overwrite [config] with [new] file")

	flag.Parse()

	var cfg *Config
	if *newConfig {
		cfg = doNewConfig(*selectedConfig, *overwriteConfig)
	} else {
		var err error
		cfg, err = core.LoadJSON[*Config](*selectedConfig)
		if err != nil {
			slog.Error("Error:%v", err)
			panic("failed to load config file")
		}
	}

	if cfg == nil {
		panic("failed to build configuration file")
	}

	appCore, err := core.New(cfg.EmrsCore)
	if err != nil {
		slog.Error("Error:%v", err)
		panic("failed to create core")
	}

	appCore.AddSnapshotReceiver(networkMapUpdateLogger)

	// TODO: Start the API server(s) and
	//       finish initializing the core

}

func networkMapUpdateLogger(nmap *core.NetworkSnapshot) {
	slog.Info("core network map updated",
		"assets", len(nmap.Assets),
		"signals", len(nmap.Signals),
		"mapped-actions", len(nmap.SignalMap))
}

func doNewConfig(name string, force bool) *Config {

	slog.Debug("making new config", "name", name, "force-overwrite-enabled", force)

	if _, err := os.Stat(name); !errors.Is(err, os.ErrNotExist) {
		if !force {
			slog.Error("Error: File [%s] alredy exists. Use --force to overwrite", name)
			os.Exit(1)
		}
	}
	cfg := CreateConfigTemplate()
	if err := cfg.WriteTo(name); err != nil {
		slog.Error("Failed to create config: %s", err.Error())
		os.Exit(1)
	}
	return cfg
}
