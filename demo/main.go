package main

import "gioui.org/ui/app"
import "gioui.org/ui"
import "gioui.org/ui/layout"
import "gophercon/simple"

func main() {
	go func() {
		theme := simple.NewTheme()
		w := app.NewWindow()
		ops := new(ui.Ops)
		for e := range w.Events() {
			switch e := e.(type) {
			case app.UpdateEvent:
				cfg := &e.Config
				ops.Reset()
				theme.Reset(cfg)
				cs := layout.RigidConstraints(e.Size)

				align := layout.Align{Alignment: layout.Center}
				cs = align.Begin(ops, cs)
				dims := theme.Label("hello, world", 46).Layout(ops, cs)
				align.End(dims)

				w.Update(ops)
			}
		}
	}()
	app.Main()
}
