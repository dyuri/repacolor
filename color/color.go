package color

import (
	"fmt"
	"image/color"

	"github.com/lucasb-eyer/go-colorful"
)

type RepaColor struct {
	colorful.Color
	A float64
}

var NoColor = RepaColor{}

const delta = 1.0 / 256.0

func almosteq_eps(a, b, eps float64) bool {
	return a-b < eps && b-a < eps
}

func almosteq(a, b float64) bool {
	return almosteq_eps(a, b, delta)
}

func formatFloat(f float64) string {
	if almosteq(f, float64(int(f))) {
		return fmt.Sprintf("%d", int(f))
	}
	return fmt.Sprintf("%.4g", f)
}

func (col RepaColor) RGBA() (r, g, b, a uint32) {
	r = uint32(col.R*0xffff + 0.5)
	g = uint32(col.G*0xffff + 0.5)
	b = uint32(col.B*0xffff + 0.5)
	a = uint32(col.A*0xffff + 0.5)
	return
}

func (col RepaColor) Hex() string {
	if col.A == 1 {
		return fmt.Sprintf("#%02x%02x%02x", uint8(col.R*255.0+0.5), uint8(col.G*255.0+0.5), uint8(col.B*255.0+0.5))
	}
	return fmt.Sprintf("#%02x%02x%02x%02x", uint8(col.R*255.0+0.5), uint8(col.G*255.0+0.5), uint8(col.B*255.0+0.5), uint8(col.A*255.0+0.5))
}

func (col RepaColor) String() string {
	return col.Hex()
}

func (col RepaColor) RgbString() string {
	if (col.A == 1) {
		return fmt.Sprintf("rgb(%d %d %d)", uint8(col.R*255.0+0.5), uint8(col.G*255.0+0.5), uint8(col.B*255.0+0.5))
	}
	return fmt.Sprintf("rgb(%d %d %d / %s)", uint8(col.R*255.0+0.5), uint8(col.G*255.0+0.5), uint8(col.B*255.0+0.5), formatFloat(col.A))
}

func (col RepaColor) HslString() string {
	h, s, l := col.Hsl()
	if (col.A == 1) {
		return fmt.Sprintf("hsl(%.3g %s%% %s%%)", h, formatFloat(s*100), formatFloat(l*100))
	}
	return fmt.Sprintf("hsl(%.3g %s%% %s%% / %s)", h, formatFloat(s*100), formatFloat(l*100), formatFloat(col.A))
}

func (col RepaColor) LabString() string {
	l, a, b := col.Lab()
	if (col.A == 1) {
		return fmt.Sprintf("lab(%.4g%% %.4g %.4g)", l*100, a*100, b*100)
	}
	return fmt.Sprintf("lab(%.4g%% %.4g %.4g / %s)", l*100, a*100, b*100, formatFloat(col.A))
}

func (col RepaColor) LchString() string {
	h, c, l := col.Hcl()
	if (col.A == 1) {
		return fmt.Sprintf("lch(%.4g%% %.4g %.4g)", l*100, c*100, h)
	}
	return fmt.Sprintf("lch(%.4g%% %.4g %.4g / %s)", l*100, c*100, h, formatFloat(col.A))
}

func (col RepaColor) OkLabString() string {
	l, a, b := col.OkLab()
	if (col.A == 1) {
		return fmt.Sprintf("oklab(%s%% %s %s)", formatFloat(l*100), formatFloat(a), formatFloat(b))
	}
	return fmt.Sprintf("oklab(%s%% %s %s / %s)", formatFloat(l*100), formatFloat(a), formatFloat(b), formatFloat(col.A))
}

func (col RepaColor) OkLchString() string {
	l, c, h := col.OkLch()
	if (col.A == 1) {
		return fmt.Sprintf("oklch(%s%% %s %s)", formatFloat(l*100), formatFloat(c), formatFloat(h))
	}
	return fmt.Sprintf("oklch(%s%% %s %s / %s)", formatFloat(l*100), formatFloat(c), formatFloat(h), formatFloat(col.A))
}

func (col RepaColor) XyzString() string {
	x, y, z := col.Xyz()
	if (col.A == 1) {
		return fmt.Sprintf("xyz(%.4g %.4g %.4g)", x, y, z)
	}
	return fmt.Sprintf("xyz(%.4g %.4g %.4g / %s)", x, y, z, formatFloat(col.A))
}

func MakeColor(col color.Color) RepaColor {
	_, _, _, a := col.RGBA()
	cc, _ := colorful.MakeColor(col)

	return RepaColor{
		cc,
		float64(a) / 0xffff,
	}
}
