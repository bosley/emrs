package main

import (
	"log/slog"

	"gioui.org/text"
	"gioui.org/widget/material"
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
)

type LoginView struct {
	changeViewFn ViewUpdateFn // changeViewFn(&myView{})

	changeViewBttn widget.Clickable
}

func NewLoginView(emrs EmrsInfo, updateFn ViewUpdateFn) View {

	view := &LoginView{
		changeViewFn: updateFn,
	}

	return view
}

func (view *LoginView) Update(cfg WindowConfig) {

	if view.changeViewBttn.Clicked(cfg.Gtx) {
		slog.Info("changing to dashboard")
		view.changeViewFn(Views["dashboard"])
		return
	}

	layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceStart,
	}.Layout(cfg.Gtx,
		layout.Rigid(
			func(ctx layout.Context) layout.Dimensions {
				title := material.H1(cfg.Theme, "LOGIN")
				colr := color.NRGBA{R: 0, G: 0, B: 144, A: 255}
				title.Color = colr

				title.Alignment = text.Middle
				return title.Layout(ctx)
			},
		),
		layout.Rigid(
			func(ctx layout.Context) layout.Dimensions {
				margins := layout.Inset{
					Top:    unit.Dp(25),
					Bottom: unit.Dp(25),
					Right:  unit.Dp(35),
					Left:   unit.Dp(35),
				}
				return margins.Layout(cfg.Gtx,
					func(ctx layout.Context) layout.Dimensions {
						btn := material.Button(cfg.Theme, &view.changeViewBttn, "to dashboard")
						return btn.Layout(ctx)
					},
				)
			},
		),
	)
}
