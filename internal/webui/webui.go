/*
  Web UI interface
*/

package webui

import (
	"context"
	"errors"
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
	DefaultWebUiAddr = "127.0.0.1:8080"
	ModeRelease      = gin.ReleaseMode
	ModeDebug        = gin.DebugMode
	DefaultWebUiMode = ModeDebug
)

const (
	webAssetDir = "internal/webui/static"
)

var ErrAlreadyStarted = errors.New("webui already started")

type Config struct {
	Engine           *nerv.Engine
	Address          string
	Mode             string
	DesignatedTopic  string
	AuthenticateUser func(username string, password string) *string
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
	pane       *nerv.ModulePane
	killOtw    atomic.Bool
	topic      string
	authUserFn func(username string, password string) *string

	metrics WebMetrics
}

func New(config Config) *WebUi {

	gin.SetMode(config.Mode)

	route := gin.New()

	ui := &WebUi{
		nrvEng:     config.Engine,
		wg:         new(sync.WaitGroup),
		running:    false,
		address:    config.Address,
		topic:      config.DesignatedTopic,
		authUserFn: config.AuthenticateUser,
	}

	ui.killOtw.Store(false)

	store := cookie.NewStore([]byte("some-badger-secret-here"))

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

func (ui *WebUi) Start() error {

	slog.Info("webui:Start")

	if err := ui.pane.SubscribeTo(ui.topic, []nerv.Consumer{
		nerv.Consumer{
			Id: ui.GetName(),
			Fn: ui.NervEvent,
		}}, true); err != nil {
		return err
	}
	if ui.running {
		return ErrAlreadyStarted
	}

	ui.srv = &http.Server{
		Addr:    ui.address,
		Handler: ui.ginEng,
	}

	ui.wg.Add(1)
	go func() {
		defer func() {
			ui.wg.Done()
			ui.running = false
		}()
		err := ui.srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			slog.Error(err.Error())
			os.Exit(-1)
		}
	}()
	ui.running = true
	return nil
}

func (ui *WebUi) Shutdown() {

	slog.Info("webui:Shutdown")

	if ui.wg == nil {
		return
	}

	shutdownCtx, shutdownRelease := context.WithTimeout(
		context.Background(), 5*time.Second)

	defer shutdownRelease()

	if err := ui.srv.Shutdown(shutdownCtx); err != nil {
		slog.Error(err.Error())
		panic("Failed to shutdown webui")
	}

	ui.wg.Wait()
	ui.wg = nil
}

func (ui *WebUi) RecvModulePane(pane *nerv.ModulePane) {
	ui.pane = pane
}

func (ui *WebUi) NervEvent(event *nerv.Event) {

	slog.Debug("webui:NervEvent")

	cmd, ok := event.Data.(*MsgCommand)
	if !ok {
		slog.Warn("webui failed to convert event to command")
		return
	}

	switch cmd.Type {
	case MsgTypeInfo:
		ui.procCmdInfo(cmd.Msg.(*MsgInfo))
		break
	default:
		slog.Warn("invalid command type for webui", "value", cmd.Type)
		break
	}
}
