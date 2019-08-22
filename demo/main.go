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
		editor := theme.Editor(46)
		for e := range w.Events() {
			switch e := e.(type) {
			case app.UpdateEvent:
				cfg := &e.Config
				ops.Reset()
				theme.Reset(cfg)
				cs := layout.RigidConstraints(e.Size)
				editor.Layout(cfg, w.Queue(), ops, cs)
				w.Update(ops)
			}
		}
	}()
	app.Main()
}
