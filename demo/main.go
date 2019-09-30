package main

import (
	"fmt"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/unit"
	"gopher.con/simple"
)

func main() {
	go func() {
		theme := simple.NewTheme()
		w := app.NewWindow()
		gtx := &layout.Context{
			Queue: w.Queue(),
		}
		list := layout.List{Axis: layout.Vertical}
		btn := new(simple.IconButton)
		n := 3
		for e := range w.Events() {
			switch e := e.(type) {
			case app.UpdateEvent:
				gtx.Reset(&e.Config, layout.RigidConstraints(e.Size))
				theme.Reset(gtx.Config)

				for btn.Clicked(gtx) {
					n += 1
				}

				list.Layout(gtx, n, func(i int) {
					s := fmt.Sprintf("hello, world %d", i)
					theme.Label(s, 46).Layout(gtx)
				})

				align := layout.Align(layout.SE)
				align.Layout(gtx, func() {
					margins := layout.UniformInset(unit.Dp(8))
					margins.Layout(gtx, func() {
						btn.Layout(gtx)
					})
				})

				w.Update(gtx.Ops)
			}
		}
	}()
	app.Main()
}
