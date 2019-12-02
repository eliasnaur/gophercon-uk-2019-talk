package main

import (
	"fmt"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

func main() {
	go func() {
		gofont.Register()
		theme := material.NewTheme()

		ico, _ := material.NewIcon(icons.ContentAdd)
		w := app.NewWindow()
		gtx := layout.NewContext(w.Queue())
		list := layout.List{Axis: layout.Vertical}
		btn := new(widget.Button)
		n := 3
		for e := range w.Events() {
			switch e := e.(type) {
			case system.FrameEvent:
				gtx.Reset(e.Config, e.Size)

				for btn.Clicked(gtx) {
					n += 1
				}

				list.Layout(gtx, n, func(i int) {
					s := fmt.Sprintf("hello, world %d", i)
					theme.H2(s).Layout(gtx)
				})

				align := layout.Align(layout.SE)
				align.Layout(gtx, func() {
					margins := layout.UniformInset(unit.Dp(8))
					margins.Layout(gtx, func() {
						b := theme.IconButton(ico)
						b.Size = unit.Dp(72)
						b.Layout(gtx, btn)
					})
				})

				e.Frame(gtx.Ops)
			}
		}
	}()
	app.Main()
}
