package main

import "gioui.org/ui/app"
import "gioui.org/ui"
import "gioui.org/ui/layout"
import "gophercon/simple"
import "fmt"

func main() {
	go func() {
		theme := simple.NewTheme()
		w := app.NewWindow()
		ops := new(ui.Ops)
		list := layout.List{Axis: layout.Vertical}
		btn := new(simple.IconButton)
		for e := range w.Events() {
			switch e := e.(type) {
			case app.UpdateEvent:
				cfg := &e.Config
				ops.Reset()
				theme.Reset(cfg)
				cs := layout.RigidConstraints(e.Size)

				for list.Init(cfg, w.Queue(), ops, cs, 100); list.More(); list.Next() {
					cs := list.Constraints()
					s := fmt.Sprintf("hello, world %d", list.Index())
					dims := theme.Label(s, 46).Layout(ops, cs)
					list.End(dims)
				}
				list.Layout()

				align := layout.Align{Alignment: layout.SE}
				cs = align.Begin(ops, cs)
				margins := layout.UniformInset(ui.Dp(8))
				cs = margins.Begin(cfg, ops, cs)
				dims := btn.Layout(cfg, ops, cs)
				dims = margins.End(dims)
				align.End(dims)

				w.Update(ops)
			}
		}
	}()
	app.Main()
}
