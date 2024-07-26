package main

import (
	"log/slog"

	"gioui.org/text"
	"gioui.org/widget/material"
	"image/color"
)

type LoginView struct {

  updateView  ViewUpdateFn    // updateView(&myView{})
}

func NewLoginView(emrs EmrsInfo, updateFn ViewUpdateFn) View {

	login := &LoginView{
    updateView: updateFn,
  }

	return login
}

func (view *LoginView) Update(cfg WindowConfig) {

	slog.Info("update window")
	// Define an large label with an appropriate text:
	title := material.H1(cfg.Theme, "Hello, Gio")

	// Change the color of the label.
	maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
	title.Color = maroon

	// Change the position of the label.
	title.Alignment = text.Middle

	// Draw the label to the graphics context.
	title.Layout(cfg.Gtx)

}
