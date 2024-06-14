/*
  Web UI interface
*/

package webui

import (
	"context"
	"errors"
	"github.com/bosley/nerv-go"
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
	DefaultWebUiMode = gin.DebugMode
)

const (
	webAssetDir = "internal/webui/static"
)

var ErrAlreadyStarted = errors.New("webui already started")

type Config struct {
	Engine  *nerv.Engine
	Address string
	Mode    string
}

type WebMetrics struct {
	requests atomic.Uint64
	// TODO:  avgLatency (see middleware.go)
}

type WebUi struct {
	ginEng    *gin.Engine
	nrvEng    *nerv.Engine
	wg        *sync.WaitGroup
	srv       *http.Server
	running   bool
	address   string
	submitter *nerv.ModuleSubmitter
	killOtw   atomic.Bool

	metrics WebMetrics
}

func New(config Config) *WebUi {

	gin.SetMode(config.Mode)

	route := gin.New()

	ui := &WebUi{
		nrvEng:  config.Engine,
		wg:      new(sync.WaitGroup),
		running: false,
		address: config.Address,
	}

	ui.killOtw.Store(false)

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
	ui.ginEng.Static("/css", strings.Join([]string{webAssetDir, "templates/css"}, "/"))
}

func (ui *WebUi) initRoutes() {
	ui.ginEng.GET("/", ui.routeHome)
	ui.ginEng.GET("/status", ui.routeStatus)
}

func (ui *WebUi) Start() error {

	slog.Info("webui:Start")

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

func (ui *WebUi) SetSubmitter(submitter *nerv.ModuleSubmitter) {
	ui.submitter = submitter
}

func (ui *WebUi) ReceiveEvent(event *nerv.Event) {

	slog.Debug("webui:ReceiveEvent")

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
