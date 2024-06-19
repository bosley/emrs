/*
  Web UI interface
*/

package webui

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/bosley/emrs/badger"
	"github.com/bosley/nerv-go"
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
	ModeRelease = gin.ReleaseMode
	ModeDebug   = gin.DebugMode
)

const (
	webAssetDir = "static"
)

var ErrAlreadyStarted = errors.New("webui already started")
var ErrNotYetStarted = errors.New("webui not yet started")

type Config struct {
	Engine           *nerv.Engine
	Address          string
	Mode             string
	KillChannel      string
	AuthenticateUser func(username string, password string) *string

	ServerCert string
	ServerKey  string
}

type WebMetrics struct {
	requests atomic.Uint64
	// TODO:  avgLatency (see middleware.go)
}

type WebUi struct {
	ginEng     *gin.Engine
	nrvEng     *nerv.Engine
	wg         *sync.WaitGroup
	srv        *http.Server
	running    bool
	address    string
	killOtw    atomic.Bool
	topic      string
	tlsConfig  *tls.Config
	authUserFn func(username string, password string) *string
	badge      badger.Badge

	metrics WebMetrics
}

func New(config Config) *WebUi {

	webuiConsumerName := "webui.consumer"

	gin.SetMode(config.Mode)

	serverTLSCert, err := tls.LoadX509KeyPair(config.ServerCert, config.ServerKey)
	if err != nil {
		slog.Error(err.Error(), "cert", config.ServerCert, "key", config.ServerKey)
		panic("failed to load cert/key")
	}

	route := gin.New()

	ui := &WebUi{
		nrvEng:     config.Engine,
		address:    config.Address,
		authUserFn: config.AuthenticateUser,
		running:    false,
		wg:         new(sync.WaitGroup),
		tlsConfig: &tls.Config{
			Certificates: []tls.Certificate{serverTLSCert},
		},
		badge: nil,
	}

	if err := setupBadger(ui); err != nil {
		// Can only happen if rand.Read fails
		slog.Error(err.Error())
		panic("failed to generate badge")
	}

	ui.killOtw.Store(false)

	// Create a nerv consumer to listen for early
	// kill signal for the reaper middleware
	ui.nrvEng.Register(nerv.Consumer{
		Id: webuiConsumerName,
		Fn: func(event *nerv.Event) {
			slog.Debug("webui received shutdown warning", "from", event.Producer)
			ui.killOtw.Store(true)
		},
	})

	if err := ui.nrvEng.SubscribeTo(config.KillChannel, webuiConsumerName); err != nil {
		panic(err.Error())
	}

	// badge id is random, so when the server reboots all
	// previous session information will be made moot
	store := cookie.NewStore([]byte(ui.badge.Id()))

	route.Use(sessions.Sessions("emrs", store))

	route.Use(ui.ReaperMiddleware())

	route.Use(ui.RequestProfiler())

	ui.ginEng = route

	ui.initRoutes()

	ui.initStatics()

	return ui
}

func (ui *WebUi) ShutdownWarning(t int) {

	ui.killOtw.Store(true)

	// We could say something about the time, but meh
}

func (ui *WebUi) initStatics() {
	ui.ginEng.LoadHTMLGlob(strings.Join([]string{webAssetDir, "templates/*.html"}, "/"))
	ui.ginEng.Static("/js", strings.Join([]string{webAssetDir, "js"}, "/emrs/"))
}

func (ui *WebUi) initRoutes() {
	ui.ginEng.GET("/", ui.routeIndex)
	ui.ginEng.GET("/login", ui.routeLogin)
	ui.ginEng.GET("/logout", ui.routeLogout)
	ui.ginEng.POST("/auth", ui.routeAuth)

	priv := ui.ginEng.Group("/emrs")
	priv.Use(ui.EmrsAuth())
	{
		priv.GET("/status", ui.routeStatus)
		priv.GET("/dashboard", ui.routeDashboard)
		priv.GET("/settings", ui.routeSettings)
	}
}

func (ui *WebUi) GetName() string {
	return "mod.webui"
}

func setupBadger(ui *WebUi) error {

	var err error
	ui.badge, err = badger.New(badger.Config{
		Nickname: "webUi",
	})

	if err != nil {
		return err
	}

	slog.Debug("badge created", "Id", ui.badge.Id(), "PubKey", ui.badge.PublicKey())
	return nil
}

func (ui *WebUi) Start() error {

	slog.Info("webui:Start")

	if ui.running {
		return ErrAlreadyStarted
	}

	ui.srv = &http.Server{
		Addr:      ui.address,
		Handler:   ui.ginEng,
		TLSConfig: ui.tlsConfig,
	}

	ui.wg.Add(1)
	go func() {
		defer func() {
			ui.wg.Done()
			ui.running = false
		}()
		err := ui.srv.ListenAndServeTLS("", "")
		if err != nil && err != http.ErrServerClosed {
			slog.Error(err.Error())
			os.Exit(-1)
		}
	}()
	ui.running = true
	return nil
}

func (ui *WebUi) Stop() error {

	slog.Info("webui:Stop")

	if !ui.running {
		return ErrNotYetStarted
	}

	shutdownCtx, shutdownRelease := context.WithTimeout(
		context.Background(), 5*time.Second)

	defer shutdownRelease()

	if err := ui.srv.Shutdown(shutdownCtx); err != nil {
		return err
	}

	ui.wg.Wait()
	return nil
}
