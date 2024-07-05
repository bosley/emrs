package main

import (
	"crypto/tls"
	"emrs/core"
	"errors"
	"flag"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
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

	/*
	   TODO:


	   Create a gin server in a public module called API eAPI or something
	   this will be what people can use to write go programs to interface with
	   the server.

	   UI Should have port and stuff removed from config. The only web thing
	   emrs should be supporting a the moment is the api and the event endpoints.

	   UI can be added later, utilizing the /api endpoint to configure/ edit the server
	   I might make the UI an entirely different app/repo

	   Public Group
	       /event    POST    Post an event to the server - Later add config to make voucher locked
	       /api
	         Middleware will grab and authenticate API vouchers for
	         every request

	         All requests are POSTS of JSON commands. Each cmmand will have the API token
	         and specific subcommands for querying/setting data


	*/
	gins := gin.New()

	gins.POST("/api", appCore.Api)

	cert, err := cfg.LoadTLSCert()

	if err != nil {
		slog.Error("Failed to load TLS Cert",
			"key", cfg.Hosting.Key,
			"crt", cfg.Hosting.Cert)
		os.Exit(1)
	}

	api := http.Server{
		Addr:    cfg.Hosting.ApiAddress,
		Handler: gins,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}

	err = api.ListenAndServeTLS("", "")
	if err != nil && err != http.ErrServerClosed {
		slog.Error(err.Error())
		os.Exit(1)
	}
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
