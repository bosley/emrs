package main

import (
	"github.com/bosley/emrs/badger"
	"github.com/bosley/emrs/datastore"

	"log/slog"
	"os"
  "time"

	//"image/color"
	"gioui.org/app"
//	"gioui.org/op"
//	"gioui.org/text"
//	"gioui.org/widget/material"
)

const (
  ViewLogin = "login"
  ViewDashboard = "dashboard"
  ViewSuite = "suite"
)

type EmrsInfo struct {
  Home string
  Cfg Config
  Ds  datastore.DataStore
  Badge badger.Badge
}

type Engine struct {
  emrs EmrsInfo
  viewQueue []string
  changeView bool

}

func MustCreateEngine(home string) *Engine {

  eng := Engine{
    viewQueue: make([]string, 1),
  }
  eng.viewQueue[0] = ViewLogin

  eng.emrs.Cfg, eng.emrs.Badge = mustLoadCfgAndBadge(home)
  eng.emrs.Ds = mustLoadDefaultDataStore(home)
  return &eng
}

func (engine *Engine) IndViewClosed() {
  slog.Info("view closed", "name", engine.viewQueue[0])
  engine.viewQueue = engine.viewQueue[1:]
}

func (engine *Engine) updateView() {
  if len(engine.viewQueue) == 0 {
    slog.Info("no more views to display")
    os.Exit(0)
  }
}

func (engine *Engine) Run() { 
  go func() {
    for {
      time.Sleep(time.Microsecond * 250)
      if engine == nil || len(engine.viewQueue) == 0 {
        slog.Debug("view-exchange thread closing")
        return
      }
      if engine.changeView {
        engine.changeView = false
        engine.updateView()
      }

      /*
    while looking at https://gioui.org/doc/architecture/window

    I think we shoudl grab the event, handle the close here,
    and then call the draw function on the view with the graphics context

      */
      switch engine.viewQueue[0] {
        case ViewLogin:

          // Each view should have a struct that is stored in a diffeent file
          // that takes in whatever to draw

          break
        case ViewDashboard:
          break
        case ViewSuite:
          break
      }
    }
  }()

  app.Main()
  slog.Info("application complete")
}


