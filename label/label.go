package main

import (
	"gioui.org/ui"
	"gioui.org/ui/app"
	"gioui.org/ui/layout"

	"gopherconuk2019/simple"
)

// START OMIT
func main() {
	go func() {
		theme := simple.NewTheme()
		ops := new(ui.Ops)
		w := app.NewWindow()
		for e := range w.Events() {
			switch e := e.(type) {
			case app.UpdateEvent:
				ops.Reset()
				theme.Reset(&e.Config)
				cs := layout.RigidConstraints(e.Size)

				theme.Label("Hello, World!", 46).Layout(ops, cs)

				w.Update(ops)
			}
		}
	}()
	app.Main()
}

// END OMIT
