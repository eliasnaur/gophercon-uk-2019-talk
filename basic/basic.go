package main

import (
	"gioui.org/ui"
	"gioui.org/ui/app"
	"gioui.org/ui/f32"
	"gioui.org/ui/paint"

	"gophercon/simple"
)

// START OMIT
func main() {
	go func() {
		_ = simple.NewTheme()
		ops := new(ui.Ops)
		w := app.NewWindow()
		for e := range w.Events() {
			switch e.(type) {
			case app.UpdateEvent:
				ops.Reset()
				//paint.ColorOp{Color: color.RGBA{A: 0xff, R: 0xff}}.Add(ops)
				paint.PaintOp{Rect: f32.Rectangle{
					Max: f32.Point{X: 500, Y: 500},
				}}.Add(ops)
				w.Update(ops)
			}
		}
	}()
	app.Main()
}

// END OMIT
