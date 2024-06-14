/*
  Web UI interface
*/

package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	defaultWebUiAddr = "127.0.0.1:8080"
	defaultWebUiMode = gin.DebugMode
)

type WebUi struct {
	engine  *gin.Engine
	wg      *sync.WaitGroup
	srv     *http.Server
	running bool
	address string
}

func CreateWebUi(mode string, address string) *WebUi {
	ui := &WebUi{
		engine:  gin.New(),
		wg:      new(sync.WaitGroup),
		running: false,
		address: address,
	}
	ui.initRoutes()
	ui.initStatics()
	return ui
}

func (ui *WebUi) Start() error {

	slog.Info("webui:Start")

	if ui.running {
		panic("webui already started")
	}

	ui.srv = &http.Server{
		Addr:    ui.address,
		Handler: ui.engine,
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

func (ui *WebUi) Stop() {

	slog.Info("webui:Stop")

	if ui.wg == nil {
		return
	}

	shutdownCtx, shutdownRelease := context.WithTimeout(
		context.Background(), 5*time.Second)

	defer shutdownRelease()

	must(ui.srv.Shutdown(shutdownCtx))

	ui.wg.Wait()
	ui.wg = nil
}

func (ui *WebUi) initStatics() {
	ui.engine.LoadHTMLGlob("web/templates/*.html")
	ui.engine.Static("/css", "web/templates/css")

}

func (ui *WebUi) initRoutes() {
	ui.engine.GET("/", ui.routeHome)
	ui.engine.GET("/status", ui.routeStatus)
}

func (ui *WebUi) routeHome(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"title":   "E.M.R.S",
		"message": "WORK IN PROGRESS",
	})
}

func (ui *WebUi) routeStatus(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "success",
	})
}
