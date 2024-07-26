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

type DashboardView struct {
	changeViewFn ViewUpdateFn // changeViewFn(&myView{})

	changeViewBttn widget.Clickable
}

func NewDashboardView(emrs EmrsInfo, updateFn ViewUpdateFn) View {

	view := &DashboardView{
		changeViewFn: updateFn,
	}

	return view
}

func (view *DashboardView) Update(cfg WindowConfig) {

	if view.changeViewBttn.Clicked(cfg.Gtx) {
		slog.Info("changing to login")
		view.changeViewFn(Views["login"])
		return
	}

	layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceStart,
	}.Layout(cfg.Gtx,
		layout.Rigid(
			func(ctx layout.Context) layout.Dimensions {
				title := material.H1(cfg.Theme, "DASHBOARD")
				colr := color.NRGBA{R: 144, G: 0, B: 0, A: 255}
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
						btn := material.Button(cfg.Theme, &view.changeViewBttn, "to login")
						return btn.Layout(ctx)
					},
				)
			},
		),
	)
}
