package webui

import (
	"context"
	"crypto/tls"
	"emrs/core"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	webAssetDir = "web-static"
)

func New(
	appCore *core.Core,
	address string,
	emrsSessionId string,
	cert tls.Certificate) *controller {
	return &controller{
		appCore: appCore,
		address: address,
		emrsId:  emrsSessionId,
		tlsCert: cert,
		wg:      new(sync.WaitGroup),
	}
}

type metricsData struct {
	requests atomic.Uint64
}

type controller struct {
	appCore *core.Core
	address string
	emrsId  string
	tlsCert tls.Certificate
	running atomic.Bool
	wg      *sync.WaitGroup
	srv     *http.Server
	killOtw atomic.Bool
	metrics metricsData
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

	store := cookie.NewStore([]byte(c.emrsId))

	gins.Use(sessions.Sessions("emrs", store))

	gins.LoadHTMLGlob(strings.Join([]string{webAssetDir, "templates/*.html"}, "/"))
	gins.Static("/js", strings.Join([]string{webAssetDir, "js"}, "/emrs/"))

	gins.GET("/", c.routeIndex)
	gins.GET("/login", c.routeLogin)
	gins.GET("/logout", c.routeLogout)
	gins.POST("/auth", c.routeAuth)

	priv := gins.Group("/emrs")
	priv.Use(c.EmrsAuth())
	{
		priv.GET("/status", c.routeStatus)
		priv.GET("/dashboard", c.routeDashboard)
		priv.GET("/settings", c.routeSettings)
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

	slog.Warn("TODO: START THE WEB SERVER")

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
