package main

import "gioui.org/ui/app"
import "gioui.org/ui/paint"
import "gioui.org/ui"
import "gioui.org/ui/f32"

func main() {
	go func() {
		w := app.NewWindow()
		ops := new(ui.Ops)
		for e := range w.Events() {
			switch e.(type) {
			case app.UpdateEvent:
				ops.Reset()
				paint.PaintOp{Rect: f32.Rectangle{
					Max: f32.Point{X: 400, Y: 500},
				}}.Add(ops)
				w.Update(ops)
			}
		}
	}()
	app.Main()
}
