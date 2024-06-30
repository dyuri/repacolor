package color

import (
	"errors"
	"fmt"
	"image/color"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
)

type RepaColor struct {
	colorful.Color
	A float64
}

var NoColor = RepaColor{}

func (col RepaColor) RGBA() (r, g, b, a uint32) {
	r = uint32(col.R * 0xffff + 0.5)
	g = uint32(col.G * 0xffff + 0.5)
	b = uint32(col.B * 0xffff + 0.5)
	a = uint32(col.A * 0xffff + 0.5)
	return
}

func (col RepaColor) Hex() string {
	if col.A == 1 {
		return fmt.Sprintf("#%02x%02x%02x", uint8(col.R*255.0+0.5), uint8(col.G*255.0+0.5), uint8(col.B*255.0+0.5))
	}
	return fmt.Sprintf("#%02x%02x%02x%02x", uint8(col.R*255.0+0.5), uint8(col.G*255.0+0.5), uint8(col.B*255.0+0.5), uint8(col.A*255.0+0.5))
}

func MakeColor(col color.Color) RepaColor {
	_, _, _, a := col.RGBA()
	cc, _ := colorful.MakeColor(col)

	return RepaColor{
		cc,
		float64(a) / 0xffff,
	}
}

func ParseColor(cstr string) (RepaColor, error) {
	// hexa color
	if strings.HasPrefix(cstr, "#") {
		if len(cstr) == 7 || len(cstr) == 4 {
			c, err := colorful.Hex(cstr)
			if err != nil {
				return NoColor, err
			}
			return RepaColor{c, 1}, nil
		} else if len(cstr) == 9 {
			c, err := colorful.Hex(cstr[:7])
			if err != nil {
				return NoColor, err
			}
			var a uint8
			n, err := fmt.Sscanf(cstr[7:], "%02x", &a)
			if err != nil || n != 1 {
				return NoColor, err
			}
			return RepaColor{c, float64(a) / 255}, nil
		} else if len(cstr) == 5 {
			c, err := colorful.Hex(cstr[:4])
			if err != nil {
				return NoColor, err
			}
			var a uint8
			n, err := fmt.Sscanf(cstr[5:], "%1x", &a)
			if err != nil || n != 1 {
				return NoColor, err
			}
			return RepaColor{c, float64(a) / 15}, nil
		}
	}
	return NoColor, errors.New("cannot parse color")
}
