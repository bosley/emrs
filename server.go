package main

import (
	"crypto/tls"
	"emrs/badger"
	"emrs/core"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"os"
)

func runServer(cfg *Config) {

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
