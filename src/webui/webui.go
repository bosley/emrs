package webui

import (
	"context"
	"crypto/tls"
	"emrs/core"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

/*
Create a new Web UI
*/
func New(
	appCore *core.Core,
	address string,
	assets string,
	cert tls.Certificate) *controller {
	return &controller{
		appCore: appCore,
		address: address,
		tlsCert: cert,
		assets:  assets,
		wg:      new(sync.WaitGroup),
	}
}

// TODO: We are collecting requests here
//
//	but we could be tracking latency and all
//	sorts of fun stuff
type metricsData struct {
	requests atomic.Uint64
}

type controller struct {

	// Directory of static assets
	assets string

	// Actual server information
	address string
	tlsCert tls.Certificate
	srv     *http.Server

	// Execution state
	running atomic.Bool
	wg      *sync.WaitGroup

	// Runtime metrics
	metrics metricsData

	// When reaping thread flags us we soft-disable the site and begin shutdown
	//  (site disabled via "ReaperMiddleware")
	killOtw atomic.Bool

	// Application core stores access to the database interface panel
	// and is how the application keeps track of state/ etc
	appCore *core.Core
}

func (c *controller) Start() error {

	if c.running.Load() {
		return nil
	}

	c.running.Store(true)

	if c.appCore.IsReleaseMode() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	gins := gin.New()

	store := memstore.NewStore(c.appCore.GetSessionKey())

	gins.Use(sessions.Sessions("emrs", store))

	gins.LoadHTMLGlob(strings.Join([]string{c.assets, "templates/*.html"}, "/"))
	gins.Static("/emrs/ui", strings.Join([]string{c.assets, "ui"}, "/"))

	gins.GET("/", c.routeIndex)
	gins.GET("/logout", c.routeLogout)
	gins.POST("/auth", c.routeAuth)

	// These endpoints are only needed the very first time the server
	// runs. Once the use has an account we don't need the endpoints.
	// They are soft-disabled once setup is complete, but this way
	// they stay off
	if c.appCore.RequiresSetup() {
		gins.POST("/new/user", c.routeNewUser)
		gins.POST("/create/user", c.routeCreateUser)
	}

	priv := gins.Group("/emrs")
	priv.Use(c.EmrsAuth())
	{
		priv.GET("/", c.routeAppLaunch)
		priv.GET("/session", c.routeSessionInfo)
	}
	c.srv = &http.Server{
		Addr:    c.address,
		Handler: gins,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{c.tlsCert},
		},
	}

	c.wg.Add(1)
	go func() {
		defer func() {
			c.wg.Done()
			c.running.Store(false)
		}()
		err := c.srv.ListenAndServeTLS("", "")
		if err != nil && err != http.ErrServerClosed {
			slog.Error(err.Error())
			os.Exit(-1)
		}
	}()

	slog.Info("webui started")
	return nil
}

func (c *controller) Stop() error {

	shutdownCtx, shutdownRelease := context.WithTimeout(
		context.Background(), 5*time.Second)

	defer shutdownRelease()

	if err := c.srv.Shutdown(shutdownCtx); err != nil {
		return err
	}

	c.wg.Wait()
	return nil
}

func (c *controller) ShutdownAlert(time.Duration) {
	slog.Warn("webui received shutdown alert")
	c.killOtw.Store(true)
}

func (c *controller) Alive() bool {
	return c.running.Load()
}
