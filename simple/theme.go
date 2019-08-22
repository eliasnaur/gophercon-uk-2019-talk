package simple

import (
	"image"
	"image/color"
	"image/draw"

	"gioui.org/ui"
	"gioui.org/ui/f32"
	"gioui.org/ui/input"
	"gioui.org/ui/layout"
	"gioui.org/ui/measure"
	"gioui.org/ui/paint"
	"gioui.org/ui/pointer"
	"gioui.org/ui/gesture"
	"gioui.org/ui/text"
	"golang.org/x/exp/shiny/iconvg"
	"golang.org/x/exp/shiny/materialdesign/icons"
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

type icon struct {
	src  []byte
	size ui.Value

	// Cached values.
	img     image.Image
	imgSize int
}

type IconButton struct {
	icon *icon
	click gesture.Click
}

func (b *IconButton) Next(queue input.Queue) (gesture.ClickEvent, bool) {
	return b.click.Next(queue)
}

func (b *IconButton) Layout(c ui.Config, ops *ui.Ops, cs layout.Constraints) layout.Dimens {
	if b.icon == nil {
		b.icon = &icon{src: icons.ContentAdd, size: ui.Dp(28)}
	}
	f := layout.Flex{Axis: layout.Vertical, Alignment: layout.End}
	f.Init(ops, cs)
	cs = f.Rigid()
	in := layout.Inset{Top: ui.Dp(4)}
	cs = in.Begin(c, ops, cs)
	col := colorMaterial(ops, rgb(0x62798c))
	dims := fab(ops, cs, b.icon.image(c), col, c.Px(ui.Dp(56)))
	pointer.EllipseAreaOp{Rect: image.Rectangle{Max: dims.Size}}.Add(ops)
	b.click.Add(ops)
	dims = in.End(dims)
	return f.Layout(f.End(dims))
}

func colorMaterial(ops *ui.Ops, color color.RGBA) ui.MacroOp {
	var mat ui.MacroOp
	mat.Record(ops)
	paint.ColorOp{Color: color}.Add(ops)
	mat.Stop()
	return mat
}

func (ic *icon) image(cfg ui.Config) image.Image {
	sz := cfg.Px(ic.size)
	if sz == ic.imgSize {
		return ic.img
	}
	m, _ := iconvg.DecodeMetadata(ic.src)
	dx, dy := m.ViewBox.AspectRatio()
	img := image.NewRGBA(image.Rectangle{Max: image.Point{X: sz, Y: int(float32(sz) * dy / dx)}})
	var ico iconvg.Rasterizer
	ico.SetDstImage(img, img.Bounds(), draw.Src)
	// Use white for icons.
	m.Palette[0] = color.RGBA{A: 0xff, R: 0xff, G: 0xff, B: 0xff}
	iconvg.Decode(&ico, ic.src, &iconvg.DecodeOptions{
		Palette: &m.Palette,
	})
	ic.img = img
	ic.imgSize = sz
	return img
}

func toRectF(r image.Rectangle) f32.Rectangle {
	return f32.Rectangle{
		Min: f32.Point{X: float32(r.Min.X), Y: float32(r.Min.Y)},
		Max: f32.Point{X: float32(r.Max.X), Y: float32(r.Max.Y)},
	}
}

func fab(ops *ui.Ops, cs layout.Constraints, ico image.Image, mat ui.MacroOp, size int) layout.Dimens {
	dp := image.Point{X: (size - ico.Bounds().Dx()) / 2, Y: (size - ico.Bounds().Dy()) / 2}
	dims := image.Point{X: size, Y: size}
	rr := float32(size) * .5
	roundRect(ops, float32(size), float32(size), rr, rr, rr, rr)
	mat.Add(ops)
	paint.PaintOp{Rect: f32.Rectangle{Max: f32.Point{X: float32(size), Y: float32(size)}}}.Add(ops)
	paint.ImageOp{Src: ico, Rect: ico.Bounds()}.Add(ops)
	paint.PaintOp{
		Rect: toRectF(ico.Bounds().Add(dp)),
	}.Add(ops)
	return layout.Dimens{Size: dims}
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
