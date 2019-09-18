package main

import (
	"fmt"

	"gioui.org/ui"
	"gioui.org/ui/app"
	"gioui.org/ui/gesture"
	"gioui.org/ui/layout"
	"gopher.con/simple"
)

func main() {
	go func() {
		theme := simple.NewTheme()
		w := app.NewWindow()
		ops := new(ui.Ops)
		list := layout.List{Axis: layout.Vertical}
		btn := new(simple.IconButton)
		n := 3
		for e := range w.Events() {
			switch e := e.(type) {
			case app.UpdateEvent:
				cfg := &e.Config
				ops.Reset()
				theme.Reset(cfg)
				cs := layout.RigidConstraints(e.Size)

				q := w.Queue()
				for e, ok := btn.Next(q); ok; e, ok = btn.Next(q) {
					if e.Type == gesture.TypeClick {
						n += 1
					}
				}

				list.Layout(cfg, q, ops, cs, n, func(cs layout.Constraints, i int) layout.Dimensions {
					s := fmt.Sprintf("hello, world %d", i)
					return theme.Label(s, 46).Layout(ops, cs)
				})

				align := layout.Align{Alignment: layout.SE}
				align.Layout(ops, cs, func(cs layout.Constraints) layout.Dimensions {
					margins := layout.UniformInset(ui.Dp(8))
					return margins.Layout(cfg, ops, cs, func(cs layout.Constraints) layout.Dimensions {
						return btn.Layout(cfg, ops, cs)
					})
				})

				w.Update(ops)
			}
		}
	}()
	app.Main()
}
