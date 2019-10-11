package main

import (
	"fmt"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/text/opentype"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"golang.org/x/image/font/gofont/goregular"
)

func main() {
	go func() {
		shaper := new(text.Shaper)
		shaper.Register(text.Font{}, opentype.Must(
			opentype.Parse(goregular.TTF),
		))
		theme := material.NewTheme(shaper)

		ico, _ := material.NewIcon(icons.ContentAdd)
		w := app.NewWindow()
		gtx := &layout.Context{
			Queue: w.Queue(),
		}
		list := layout.List{Axis: layout.Vertical}
		btn := new(widget.Button)
		n := 3
		for e := range w.Events() {
			switch e := e.(type) {
			case app.FrameEvent:
				gtx.Reset(&e.Config, e.Size)

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
