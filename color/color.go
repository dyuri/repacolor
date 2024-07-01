package color

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/lucasb-eyer/go-colorful"
)

type RepaColor struct {
	colorful.Color
	A float64
}

var NoColor = RepaColor{}

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
	return fmt.Sprintf("rgb(%d %d %d / %s)", uint8(col.R*255.0+0.5), uint8(col.G*255.0+0.5), uint8(col.B*255.0+0.5), strconv.FormatFloat(col.A, 'f', -1, 64))
}

func MakeColor(col color.Color) RepaColor {
	_, _, _, a := col.RGBA()
	cc, _ := colorful.MakeColor(col)

	return RepaColor{
		cc,
		float64(a) / 0xffff,
	}
}
