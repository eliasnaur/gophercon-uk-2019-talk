package simple

import (
	"image"
	"image/color"
	"image/draw"
	"time"

	"gioui.org/f32"
	"gioui.org/gesture"
	"gioui.org/layout"
	"gioui.org/measure"
	"gioui.org/paint"
	"gioui.org/pointer"
	"gioui.org/text"
	"gioui.org/ui"
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
	icon     *icon
	click    gesture.Click
	clicks   int
	inkStart bool
	inkTime  time.Time
	inkPos   f32.Point
}

func (b *IconButton) Clicked(gtx *layout.Context) bool {
	for _, e := range b.click.Events(gtx) {
		if e.Type == gesture.TypeClick {
			b.inkPos = e.Position
			b.inkStart = true
			b.clicks++
		}
	}
	if b.clicks > 0 {
		b.clicks--
		if b.clicks > 0 {
			ui.InvalidateOp{}.Add(gtx.Ops)
		}
		return true
	}
	return false
}

func (b *IconButton) Layout(gtx *layout.Context) {
	if b.icon == nil {
		b.icon = &icon{src: icons.ContentAdd, size: ui.Dp(48)}
	}
	f := layout.Flex{Axis: layout.Vertical, Alignment: layout.End}
	f.Init(gtx)
	child := f.Rigid(func() {
		in := layout.Inset{Top: ui.Dp(8)}
		in.Layout(gtx, func() {
			col := colorMaterial(gtx.Ops, rgb(0x00dd00))
			if b.click.State() == gesture.StatePressed {
				col = colorMaterial(gtx.Ops, rgb(0x00aa00))
			}
			size := gtx.Px(ui.Dp(112))
			fab(gtx, b.icon.image(gtx), col, size)
			if b.inkStart {
				b.inkTime = gtx.Now()
				b.inkStart = false
			}
			if d := gtx.Now().Sub(b.inkTime); d < time.Second {
				t := float32(d.Seconds())
				var stack ui.StackOp
				stack.Push(gtx.Ops)
				size := float32(size) * 7 * t
				rr := size * .5
				col := byte(0xaa * (1 - t*t))
				ink := colorMaterial(gtx.Ops, color.RGBA{A: col, R: col, G: col, B: col})
				ink.Add(gtx.Ops)
				ui.TransformOp{}.Offset(b.inkPos).Offset(f32.Point{
					X: -rr,
					Y: -rr,
				}).Add(gtx.Ops)
				roundRect(gtx.Ops, float32(size), float32(size), rr, rr, rr, rr)
				paint.PaintOp{Rect: f32.Rectangle{Max: f32.Point{X: float32(size), Y: float32(size)}}}.Add(gtx.Ops)
				stack.Pop()
				ui.InvalidateOp{}.Add(gtx.Ops)
			}
			pointer.EllipseAreaOp{Rect: image.Rectangle{Max: gtx.Dimensions.Size}}.Add(gtx.Ops)
			b.click.Add(gtx.Ops)
		})
	})
	f.Layout(child)
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

func fab(gtx *layout.Context, ico image.Image, mat ui.MacroOp, size int) {
	dp := image.Point{X: (size - ico.Bounds().Dx()) / 2, Y: (size - ico.Bounds().Dy()) / 2}
	dims := image.Point{X: size, Y: size}
	rr := float32(size) * .5
	roundRect(gtx.Ops, float32(size), float32(size), rr, rr, rr, rr)
	mat.Add(gtx.Ops)
	paint.PaintOp{Rect: f32.Rectangle{Max: f32.Point{X: float32(size), Y: float32(size)}}}.Add(gtx.Ops)
	paint.ImageOp{Src: ico, Rect: ico.Bounds()}.Add(gtx.Ops)
	paint.PaintOp{
		Rect: toRectF(ico.Bounds().Add(dp)),
	}.Add(gtx.Ops)
	gtx.Dimensions = layout.Dimensions{Size: dims}
}

func (r Rect) Layout(gtx *layout.Context) {
	cs := gtx.Constraints
	col := argb(r.Color)
	sz := image.Point{X: cs.Width.Max, Y: cs.Height.Max}
	if rr := r.Corner; rr > 0 {
		roundRect(gtx.Ops, float32(sz.X), float32(sz.Y), rr, rr, rr, rr)
	}
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: f32.Rectangle{
		Max: f32.Point{
			X: float32(sz.X),
			Y: float32(sz.Y),
		},
	}}.Add(gtx.Ops)
	gtx.Dimensions = layout.Dimensions{Size: sz}
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

const Blah = `1. I learned from my grandfather, Verus, to use good manners, and to
put restraint on anger. 2. In the famous memory of my father I had a
pattern of modesty and manliness. 3. Of my mother I learned to be
pious and generous; to keep myself not only from evil deeds, but even
from evil thoughts; and to live with a simplicity which is far from
customary among the rich. 4. I owe it to my great-grandfather that I
did not attend public lectures and discussions, but had good and able
teachers at home; and I owe him also the knowledge that for things of
this nature a man should count no expense too great.

5. My tutor taught me not to favour either green or blue at the
chariot races, nor, in the contests of gladiators, to be a supporter
either of light or heavy armed. He taught me also to endure labour;
not to need many things; to serve myself without troubling others; not
to intermeddle in the affairs of others, and not easily to listen to
slanders against them.

6. Of Diognetus I had the lesson not to busy myself about vain things;
not to credit the great professions of such as pretend to work
wonders, or of sorcerers about their charms, and their expelling of
Demons and the like; not to keep quails (for fighting or divination),
nor to run after such things; to suffer freedom of speech in others,
and to apply myself heartily to philosophy. Him also I must thank for
my hearing first Bacchius, then Tandasis and Marcianus; that I wrote
dialogues in my youth, and took a liking to the philosopher's pallet
and skins, and to the other things which, by the Grecian discipline,
belong to that profession.

7. To Rusticus I owe my first apprehensions that my nature needed
reform and cure; and that I did not fall into the ambition of the
common Sophists, either by composing speculative writings or by
declaiming harangues of exhortation in public; further, that I never
strove to be admired by ostentation of great patience in an ascetic
life, or by display of activity and application; that I gave over the
study of rhetoric, poetry, and the graces of language; and that I did
not pace my house in my senatorial robes, or practise any similar
affectation. I observed also the simplicity of style in his letters,
particularly in that which he wrote to my mother from Sinuessa. I
learned from him to be easily appeased, and to be readily reconciled
with those who had displeased me or given cause of offence, so soon as
they inclined to make their peace; to read with care; not to rest
satisfied with a slight and superficial knowledge; nor quickly to
assent to great talkers. I have him to thank that I met with the
discourses of Epictetus, which he furnished me from his own library.

8. From Apollonius I learned true liberty, and tenacity of purpose; to
regard nothing else, even in the smallest degree, but reason always;
and always to remain unaltered in the agonies of pain, in the losses
of children, or in long diseases. He afforded me a living example of
how the same man can, upon occasion, be most yielding and most
inflexible. He was patient in exposition; and, as might well be seen,
esteemed his fine skill and ability in teaching others the principles
of philosophy as the least of his endowments. It was from him that I
learned how to receive from friends what are thought favours without
seeming humbled by the giver or insensible to the gift.`
