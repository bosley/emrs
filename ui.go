package main

import (
	"crypto/tls"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"os"
)

func runUi(cfg *Config) {
	gins := gin.New()
	gins.LoadHTMLGlob("web/templates/*.html")
	gins.Static("/img", "web/img/")
	gins.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ready",
		})
	})

	cert, err := cfg.LoadTLSCert()

	if err != nil {
		slog.Error("Failed to load TLS Cert",
			"key", cfg.Hosting.Key,
			"crt", cfg.Hosting.Cert)
		os.Exit(1)
	}

	ui := http.Server{
		Addr:    cfg.Hosting.UiAddress,
		Handler: gins,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}

	err = ui.ListenAndServeTLS("", "")
	if err != nil && err != http.ErrServerClosed {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
