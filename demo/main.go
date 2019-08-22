package main

import "gioui.org/ui/app"

func main() {
	go func() {
		w := app.NewWindow()
		for range w.Events() {
		}
	}()
	app.Main()
}
