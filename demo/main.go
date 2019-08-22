package main

import "gioui.org/ui/app"
import "gioui.org/ui"

func main() {
	go func() {
		w := app.NewWindow()
		ops := new(ui.Ops)
		for e := range w.Events() {
			switch e.(type) {
			case app.UpdateEvent:
				w.Update(ops)
			}
		}
	}()
	app.Main()
}
