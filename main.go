package main

import (
	"emrs/core"
	"errors"
	"flag"
	"github.com/gin-gonic/gin"
	"log/slog"
	"os"
)

const (
	emrsVersion            = "0.0.0"
	emrsActionScriptPrefix = "action_"
)

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
	uiEnabled := flag.Bool("ui", true, "Enable the UI endpoint for configuring the server setup")

	flag.Parse()

	cfg := loadConfig(*selectedConfig, *newConfig, *overwriteConfig)

	if cfg.Runtime.Mode == "rel" ||
		cfg.Runtime.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	slog.Info("Launching EMRS Server")
	runServer(cfg, *uiEnabled)
	return
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

func loadConfig(path string, isNew bool, force bool) *Config {
	var cfg *Config
	if isNew {
		cfg = doNewConfig(path, force)
	} else {
		var err error
		cfg, err = core.LoadJSON[*Config](path)
		if err != nil {
			slog.Error("Error:%v", err)
			panic("failed to load config file")
		}
	}

	if cfg == nil {
		panic("failed to build configuration file")
	}

	if e := cfg.Validate(); e != nil {
		slog.Error("Error:%v", e)
		panic("invalid configuration")
	}

	return cfg.WithSavePath(path)
}
