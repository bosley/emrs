package main

import (
	"crypto/tls"
	"emrs/badger"
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

	if e := cfg.Validate(); e != nil {
		slog.Error("Error:%v", e)
		panic("invalid configuration")
	}

	appCore, err := core.New(cfg.EmrsCore)
	if err != nil {
		slog.Error("Error:%v", err)
		panic("failed to create core")
	}

	appCore.AddSnapshotReceiver(func(ns *core.NetworkSnapshot) {
		slog.Info("core network map updated",
			"assets", len(ns.Assets),
			"signals", len(ns.Signals),
			"mapped-actions", len(ns.SignalMap))
	})

	gins := gin.New()

	if cfg.Runtime.Mode == "rel" ||
		cfg.Runtime.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	gins.POST("/", buildSubmit(appCore))
	gins.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ready",
		})
	})

	priv := gins.Group("/api")
	priv.Use(buildApiAuthMiddleware(
		appCore.GetPublicKey(),
		cfg.Hosting.ApiKeys,
	))

	{
		priv.GET("/", buildApi(appCore))
		priv.POST("/update", buildApiUpdate(appCore))
	}

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

func buildSubmit(app *core.Core) func(*gin.Context) {
	return func(c *gin.Context) {

		type EventSubmit struct {
			Origin string `json:"origin"`
			Data   string `json:"data"`
		}

		var es EventSubmit
		c.BindJSON(&es)

		if err := app.SubmitEvent(es.Origin, es.Data); err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"status": "okay",
		})
	}
}

func buildApiAuthMiddleware(pk string, tokens []string) func(*gin.Context) {
	authSet := core.SetFrom(tokens)

	return func(c *gin.Context) {
		key, ok := c.GetQuery("key")
		if !ok {
			slog.Error("no api key present")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "No api key given",
			})
			c.Abort()
		}

		if !authSet.Contains(key) {
			slog.Error("key not in known set of vouchers", "key", key)
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "Invalid api key",
			})
			c.Abort()
		}

		if !badger.ValidateVoucher(pk, key) {
			slog.Error("badger failed to validate key", "key", key)
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "Invalid api key",
			})
			c.Abort()
		}
	}
}

func buildApi(app *core.Core) func(*gin.Context) {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"public_key": app.GetPublicKey(),
			"topo":       app.GetTopo(),
		})
	}
}

func buildApiUpdate(app *core.Core) func(*gin.Context) {
	return func(c *gin.Context) {

		slog.Debug("update topo request")

		type UpdateSubmission struct {
			NewTopo core.Topo `json:"topo"`
		}

		var us UpdateSubmission
		c.BindJSON(&us)

		if err := app.UpdateNetworkMap(us.NewTopo); err != nil {
			slog.Error("failed to update topo", "error", err.Error())
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		slog.Info("network map updated")
		c.JSON(200, gin.H{
			"status": "okay",
		})
	}
}
