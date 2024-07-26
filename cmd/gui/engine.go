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

type Engine struct {
	UserData interface{}

	viewQueue  []View
	changeView bool
	window     *app.Window
}

func MustCreateEngine() *Engine {
	eng := Engine{
		viewQueue: make([]View, 0),
		window:    new(app.Window),
	}
	return &eng
}

func (engine *Engine) PushView(v View) *Engine {
	engine.viewQueue = append(engine.viewQueue, v)
	return engine
}

func (engine *Engine) IndViewClosed() {
	slog.Info("view closed", "name", engine.viewQueue[0])
	engine.viewQueue = engine.viewQueue[1:]
	if len(engine.viewQueue) == 0 {
		slog.Info("complete")
		os.Exit(0)
	}
}

func (engine *Engine) Run() {
	go func() {
		theme := material.NewTheme()
		var ops op.Ops

		for {
			if engine == nil || len(engine.viewQueue) == 0 {
				slog.Debug("view-exchange thread closing")
				return
			}

			switch event := engine.window.Event().(type) {
			case app.DestroyEvent:
				engine.IndViewClosed()
				break
			case app.FrameEvent:
				cfg := WindowConfig{
					Gtx:   app.NewContext(&ops, event),
					Theme: theme,
				}
				engine.viewQueue[0].Update(cfg)
				event.Frame(cfg.Gtx.Ops)
				break
			}
		}
	}()

	app.Main()
}
