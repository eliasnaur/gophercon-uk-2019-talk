package simple

import (
	"image"
	"image/color"

	"gioui.org/ui"
	"gioui.org/ui/f32"
	"gioui.org/ui/layout"
	"gioui.org/ui/measure"
	"gioui.org/ui/paint"
	"gioui.org/ui/text"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/sfnt"
)

type Theme struct {
	faces   measure.Faces
	regular *sfnt.Font
}

type Rect struct {
	// ARGB color.
	Color uint32
	// Corner radius for round rects.
	Corner float32
}

func NewTheme() *Theme {
	regular, err := sfnt.Parse(goregular.TTF)
	if err != nil {
		// Parsing Go fonts should never fail.
		panic(err)
	}
	return &Theme{
		regular: regular,
	}
}

func (t *Theme) Reset(c ui.Config) {
	t.faces.Reset(c)
}

func (t *Theme) face(size float32) text.Face {
	return t.faces.For(t.regular, ui.Sp(size))
}

func (t *Theme) Editor(size float32) *text.Editor {
	return &text.Editor{
		Face: t.face(size),
	}
}

func (t *Theme) Label(txt string, size float32) text.Label {
	return text.Label{Face: t.face(size), Text: txt}
}

func (r Rect) Layout(ops *ui.Ops, cs layout.Constraints) layout.Dimens {
	col := argb(r.Color)
	sz := image.Point{X: cs.Width.Max, Y: cs.Height.Max}
	if rr := r.Corner; rr > 0 {
		roundRect(ops, float32(sz.X), float32(sz.Y), rr, rr, rr, rr)
	}
	paint.ColorOp{Color: col}.Add(ops)
	paint.PaintOp{Rect: f32.Rectangle{
		Max: f32.Point{
			X: float32(sz.X),
			Y: float32(sz.Y),
		},
	}}.Add(ops)
	return layout.Dimens{Size: sz}
}

func rgb(c uint32) color.RGBA {
	return argb((0xff << 24) | c)
}

func argb(c uint32) color.RGBA {
	return color.RGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}

// START RR OMIT
// https://pomax.github.io/bezierinfo/#circles_cubic.
func roundRect(ops *ui.Ops, width, height, se, sw, nw, ne float32) {
	w, h := float32(width), float32(height)
	const c = 0.55228475 // 4*(sqrt(2)-1)/3
	var b paint.PathBuilder
	b.Init(ops)
	b.Move(f32.Point{X: w, Y: h - se})
	b.Cube(f32.Point{X: 0, Y: se * c}, f32.Point{X: -se + se*c, Y: se}, f32.Point{X: -se, Y: se})
	b.Line(f32.Point{X: sw - w + se, Y: 0})
	b.Cube(f32.Point{X: -sw * c, Y: 0}, f32.Point{X: -sw, Y: -sw + sw*c}, f32.Point{X: -sw, Y: -sw})
	b.Line(f32.Point{X: 0, Y: nw - h + sw})
	b.Cube(f32.Point{X: 0, Y: -nw * c}, f32.Point{X: nw - nw*c, Y: -nw}, f32.Point{X: nw, Y: -nw})
	b.Line(f32.Point{X: w - ne - nw, Y: 0})
	b.Cube(f32.Point{X: ne * c, Y: 0}, f32.Point{X: ne, Y: ne - ne*c}, f32.Point{X: ne, Y: ne})
	b.End()
}
