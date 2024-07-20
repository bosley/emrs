package app

import (
	"errors"
	"github.com/bosley/emrs/badger"
	"github.com/bosley/emrs/datastore"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

type Opts struct {
	Badge     badger.Badge
	Binding   string
	DataStore datastore.DataStore
}

type httpsInfo struct {
	keyPath  string
	certPath string
}

type App struct {
	binding string
	badge   badger.Badge
	db      datastore.DataStore
	started time.Time

	httpsSettings *httpsInfo // nil if not using https
}

func New(options *Opts) *App {
	return &App{
		binding: options.Binding,
		badge:   options.Badge,
		db:      options.DataStore,
	}
}

func (a *App) UseHttps(keyPath string, certPath string) {
	a.httpsSettings = &httpsInfo{
		keyPath:  keyPath,
		certPath: certPath,
	}
}

func (a *App) Run(enableReleaseMode bool) {

	slog.Info("Application starting", "binding", a.binding)
	a.started = time.Now()

	if enableReleaseMode {
		slog.Info("Release mode enabled")
		gin.SetMode(gin.ReleaseMode)
	}

	gins := gin.New()

	// Command and control API
	//
	//        /cnc
	//
	// Every request must contain the uikey from the
	// configuration file, and that key must be a
	// valid badger voucher created from the server's identity
	a.setupCNC(gins)

	// Public facing statistics
	//
	//      /stat
	//
	// Get statistics stuff. Metrics, etc. The `/` path
	// will be how external apis will tell if the server is up
	a.setupStat(gins)

	// Public facing submissions
	//
	//      /submit
	//
	// Submission of events, and eventually, submission of
	// information (gossip/etc) from other emrs instances
	a.setupSubmit(gins)

	var err error
	if a.httpsSettings != nil {
		slog.Info("Using TLS")
		err = gins.RunTLS(
			a.binding,
			a.httpsSettings.certPath,
			a.httpsSettings.keyPath)
	} else {
		slog.Warn("Not using TLS")
		err = gins.Run(a.binding)
	}

	if errors.Is(err, http.ErrServerClosed) {
	} else if err != nil {
		slog.Error("error starting the server", "error", err.Error())
		os.Exit(1)
	}
}

// All HTTP requests come with two pieces of information to validate them
// and permit the request:
//
//	origin:     The Asset id of the thing submitting data that must
//	            be known by the server
//	token:      A badger voucher that must be valid
func (a *App) validateRequest(origin string, token string) error {

	slog.Debug("validate request", "origin", origin, "token", token)

	if strings.TrimSpace(origin) == "" {
		return errors.New("invalid origin data")
	}

	if strings.TrimSpace(token) == "" {
		return errors.New("invalid token data")
	}

	if !a.db.AssetExists(origin) {
		slog.Error("unknown originating asset given in header", "origin", origin)
		return errors.New("unknown asset")
	}

	if !a.badge.ValidateVoucher(token) {
		return errors.New("invalid token")
	}

	return nil
}
