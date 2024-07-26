package main

import (
	"os"
	"log/slog"
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
)

type WindowConfig struct {
	Gtx   layout.Context
	Theme *material.Theme
}

type View interface {
	Update(cfg WindowConfig)
}

var Views map[string]View

type ViewUpdateFn func(view View)

type Engine struct {
  activeView View
	window     *app.Window
}

func MustCreateEngine() *Engine {
	eng := Engine{
    activeView: nil,
		window:    new(app.Window),
	}
	return &eng
}

func (engine *Engine) GetViewUpdateFn() ViewUpdateFn {
  return func(view View) {
    engine.activeView = view
  }
}

func (engine *Engine) Run() {
	go func() {
		theme := material.NewTheme()
		var ops op.Ops
		for {
			if engine == nil || engine.activeView == nil {
				slog.Debug("view-exchange thread closing")
				return
			}

			switch event := engine.window.Event().(type) {
			case app.DestroyEvent:
        slog.Info("close event")
        os.Exit(0)
			case app.FrameEvent:
				cfg := WindowConfig{
					Gtx:   app.NewContext(&ops, event),
					Theme: theme,
				}
        if engine.activeView == nil {
          slog.Error("no active view")
          continue
        }
	      engine.activeView.Update(cfg)
				event.Frame(cfg.Gtx.Ops)
			}
		}
	}()

	app.Main()
}
