package color

import (
	"fmt"
	"image/color"
	"math"

	"github.com/lucasb-eyer/go-colorful"
)

type RepaColor struct {
	colorful.Color
	A float64
}

var NOCOLOR = RepaColor{}
var BLACK = RepaColor{colorful.Color{R: 0, G: 0, B: 0}, 1}
var WHITE = RepaColor{colorful.Color{R: 1, G: 1, B: 1}, 1}
var GRAY = RepaColor{colorful.Color{R: .5, G: .5, B: .5}, 1}
var DARKGRAY = RepaColor{colorful.Color{R: .25, G: .25, B: .25}, 1}
var LIGHTGRAY = RepaColor{colorful.Color{R: .75, G: .75, B: .75}, 1}

var ANSI_RESET = "\033[0m"

const Delta = colorful.Delta

func almosteq_eps(a, b, eps float64) bool {
	return a-b < eps && b-a < eps
}

func almosteq(a, b float64) bool {
	return almosteq_eps(a, b, Delta)
}

func formatFloat(f float64) string {
	if almosteq(f, float64(int(f))) {
		return fmt.Sprintf("%d", int(f))
	}
	return fmt.Sprintf("%.4g", f)
}

func CreateColor(mode string, v1, v2, v3, a float64) RepaColor {
	switch mode {
	case "rgb":
		return RepaColor{colorful.Color{R: v1, G: v2, B: v3}, a}
	case "hsl":
		return RepaColor{colorful.Hsl(v1, v2, v3), a}
	case "lab":
		return RepaColor{colorful.Lab(v1, v2, v3), a}
	case "lch":
		return RepaColor{colorful.Hcl(v3, v2, v1), a}
	case "hcl":
		return RepaColor{colorful.Hcl(v1, v2, v3), a}
	case "oklab":
		return RepaColor{colorful.OkLab(v1, v2, v3), a}
	case "oklch":
		return RepaColor{colorful.OkLch(v1, v2, v3), a}
	case "xyz":
		return RepaColor{colorful.Xyz(v1, v2, v3), a}
	}

	// use rgb as fallback
	return RepaColor{colorful.Color{R: v1, G: v2, B: v3}, a}
}

func (col RepaColor) RGBA() (r, g, b, a uint32) {
	r = uint32(col.A*col.R*0xffff + 0.5)
	g = uint32(col.A*col.G*0xffff + 0.5)
	b = uint32(col.A*col.B*0xffff + 0.5)
	a = uint32(col.A*0xffff + 0.5)
	return
}

func (col RepaColor) RGB256() (r, g, b uint8) {
	r = uint8(col.R*255.0 + 0.5)
	g = uint8(col.G*255.0 + 0.5)
	b = uint8(col.B*255.0 + 0.5)
	return
}

func (col RepaColor) RGBA256() (r, g, b, a uint8) {
	r = uint8(col.R*255.0 + 0.5)
	g = uint8(col.G*255.0 + 0.5)
	b = uint8(col.B*255.0 + 0.5)
	a = uint8(col.A*255.0 + 0.5)
	return
}

func (col RepaColor) Hex() string {
	r, g, b, a := col.RGBA256()
	if col.A == 1 {
		return fmt.Sprintf("#%02x%02x%02x", r, g, b)
	}
	return fmt.Sprintf("#%02x%02x%02x%02x", r, g, b, a)
}

func (col RepaColor) String() string {
	return col.Hex()
}

func (col RepaColor) RgbString() string {
	r, g, b := col.RGB256()
	if (col.A == 1) {
		return fmt.Sprintf("rgb(%d %d %d)", r, g, b)
	}
	return fmt.Sprintf("rgb(%d %d %d / %s)", r, g, b, formatFloat(col.A))
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
		return fmt.Sprintf("lab(%s%% %s %s)", formatFloat(l*100), formatFloat(a*100), formatFloat(b*100))
	}
	return fmt.Sprintf("lab(%s%% %s %s / %s)", formatFloat(l*100), formatFloat(a*100), formatFloat(b*100), formatFloat(col.A))
}

func (col RepaColor) LchString() string {
	h, c, l := col.Hcl()
	if (col.A == 1) {
		return fmt.Sprintf("lch(%s%% %s %s)", formatFloat(l*100), formatFloat(c*100), formatFloat(h))
	}
	return fmt.Sprintf("lch(%s%% %s %s / %s)", formatFloat(l*100), formatFloat(c*100), formatFloat(h), formatFloat(col.A))
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

func (col RepaColor) Luminance() float64 {
	return 0.21263900587151036*col.R + 0.71516867876775593*col.G + 0.072192315360733715*col.B
}

func (col RepaColor) ContrastRatio(c2 RepaColor) float64 {
	l1 := col.Luminance()
	l2 := c2.Luminance()
	if l1 > l2 {
		return (l1 + 0.05) / (l2 + 0.05)
	}
	return (l2 + 0.05) / (l1 + 0.05)
}

func (col RepaColor) A11YPair() RepaColor {
	// (x + .05) / 0.05 = 1.05 / (x + .05) => 0.179
	if col.Luminance() > 0.179 {
		return BLACK
	}
	return WHITE
}

func (col RepaColor) AnsiFg() string {
	r, g, b := col.RGB256()
	return fmt.Sprintf("\033[38;2;%d;%d;%d;1m", r, g, b)
}

func (col RepaColor) AnsiBg() string {
	r, g, b := col.RGB256()
	pr, pg, pb := col.A11YPair().RGB256()
	return fmt.Sprintf("\033[48;2;%d;%d;%d;38;2;%d;%d;%d;1m", r, g, b, pr, pg, pb)
}

// Blend two colors based on their alpha value 
// `gamma` is the gamma correction value, default is 2.2
func (col RepaColor) AlphaBlendRgb(c2 RepaColor, gamma float64) RepaColor {
	a1 := col.A
	a2 := c2.A
	a := a1 + a2*(1-a1)
	if a == 0 {
		return NOCOLOR
	}

	if gamma == 0 {
		gamma = 2.2
	}

	return RepaColor{
		colorful.Color{
			R: math.Pow(math.Pow(col.R, gamma)*a1 + math.Pow(c2.R, gamma)*a2*(1-a1), 1/gamma),
			G: math.Pow(math.Pow(col.G, gamma)*a1 + math.Pow(c2.G, gamma)*a2*(1-a1), 1/gamma),
			B: math.Pow(math.Pow(col.B, gamma)*a1 + math.Pow(c2.B, gamma)*a2*(1-a1), 1/gamma),
		},
		a,
	}
}

func MakeColor(col color.Color) RepaColor {
	_, _, _, a := col.RGBA()
	cc, _ := colorful.MakeColor(col)

	return RepaColor{
		cc,
		float64(a) / 0xffff,
	}
}
