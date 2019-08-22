package main

import "gophercon/simple"

import "gioui.org/ui"
import "gioui.org/ui/app"
import "gioui.org/ui/layout"
import "log"

func main() {
	go func() {
		theme := simple.NewTheme()
		w := app.NewWindow()
		ops := new(ui.Ops)
		btn := new(simple.IconButton)
		q := w.Queue()
		for e := range w.Events() {
			switch e := e.(type) {
			case app.UpdateEvent:
				cfg := &e.Config
				ops.Reset()
				theme.Reset(cfg)
				cs := layout.RigidConstraints(e.Size)
				for e, ok := btn.Next(q); ok; e, ok = btn.Next(q) {
					log.Println(e)
				}
				btn.Layout(cfg, ops, cs)
				w.Update(ops)
			}
		}
	}()
	app.Main()
}
