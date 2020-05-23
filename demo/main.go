package main

import (
	"fmt"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

func main() {
	go func() {
		gofont.Register()
		theme := material.NewTheme()

		ico, _ := widget.NewIcon(icons.ContentAdd)
		w := app.NewWindow()
		var ops op.Ops
		list := layout.List{Axis: layout.Vertical}
		btn := new(widget.Clickable)
		n := 3
		for e := range w.Events() {
			switch e := e.(type) {
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e.Queue, e.Config, e.Size)

				for btn.Clicked(gtx) {
					n += 1
				}

				list.Layout(gtx, n, func(gtx layout.Context, i int) layout.Dimensions {
					s := fmt.Sprintf("hello, world %d", i)
					return material.H2(theme, s).Layout(gtx)
				})

				layout.SE.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					margins := layout.UniformInset(unit.Dp(8))
					return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						b := material.IconButton(theme, btn, ico)
						b.Size = unit.Dp(72)
						return b.Layout(gtx)
					})
				})

				e.Frame(gtx.Ops)
			}
		}
	}()
	app.Main()
}
