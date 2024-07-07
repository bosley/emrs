package main

import (
	"crypto/tls"
	"emrs/badger"
	"emrs/core"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"os"
)

func runServer(cfg *Config, uiEnabled bool) {

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
	priv := gins.Group("/api")
	priv.Use(buildApiAuthMiddleware(
		appCore.GetPublicKey(),
		cfg.Hosting.ApiKeys,
	))

	{
		priv.GET("/", buildApi(appCore))
		priv.POST("/update", buildApiUpdate(appCore))
		priv.GET("/topo", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"topo": appCore.GetTopo(),
			})
		})
	}

	if uiEnabled {
		if len(cfg.Hosting.ApiKeys) == 0 {
			panic("No available API keys to use for UI interaction")
		}
		gins.LoadHTMLGlob("web/templates/*.html")
		gins.Static("/img", "web/img/")
		gins.Static("/app", "web/app/")
		gins.GET("/", func(c *gin.Context) {
			c.HTML(200, "index.html", gin.H{
				"KeyParam": fmt.Sprintf(
					"?key=%s",
					cfg.Hosting.ApiKeys[0]),
			})
		})

	} else {
		gins.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ready",
			})
		})
	}

	cert, err := cfg.LoadTLSCert()

	if err != nil {
		slog.Error("Failed to load TLS Cert",
			"key", cfg.Hosting.Key,
			"crt", cfg.Hosting.Cert)
		os.Exit(1)
	}

	api := http.Server{
		Addr:    cfg.Hosting.Address,
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

		slog.Debug("api update")

		var msg ApiMsg
		c.BindJSON(&msg)

		currentMap, err := core.NetworkMapFromTopo(app.GetRawTopo())

		if err != nil {
			c.JSON(500, gin.H{
				"reason": "Unable to load current network map",
			})
		}

		validOps := core.SetFrom([]string{
			OpAdd,
			OpDel,
		})

		validSubjects := core.SetFrom([]string{
			SubjectSector,
			SubjectAsset,
			SubjectSignal,
			SubjectAction,
			SubjectMapping,
			SubjectTopo,
		})

		if !validSubjects.Contains(msg.Subject) {
			slog.Error("invalid subject", "data", msg.Subject)
			c.JSON(400, gin.H{
				"reason": "Invalid subject",
			})
			return
		}

		if !validOps.Contains(msg.Op) {
			slog.Error("invalid op", "data", msg.Op)
			c.JSON(400, gin.H{
				"reason": "Invalid operation",
			})
			return
		}

		err = nil

		switch msg.Subject {
		case SubjectSector:
			switch msg.Op {
			case OpAdd:
				x := new(core.Sector)
				if e := json.Unmarshal([]byte(msg.Data), x); e != nil {
					c.JSON(400, gin.H{
						"reason": "Failed to decode message data",
					})
					return
				}
				if e := currentMap.AddSector(x); e != nil {
					slog.Error(e.Error())
					c.JSON(400, gin.H{
						"reason": e.Error(),
					})
					return
				}
				break
			case OpDel:
				currentMap.DeleteSector(msg.Data)
				break
			}
			break
		case SubjectAsset:
		case SubjectSignal:
		case SubjectAction:
		case SubjectMapping:
		case SubjectTopo:
		}

		topo := currentMap.ToTopo()

		if e := app.UpdateNetworkMap(topo); e != nil {
			slog.Error(e.Error())
			c.JSON(400, gin.H{
				"reason": e.Error(),
			})
			return
		}
		c.JSON(200, gin.H{"result": "complete"})
		slog.Debug("updated", "topo", app.GetTopo())
	}
}
