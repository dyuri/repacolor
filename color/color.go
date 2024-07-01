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

func MakeColor(col color.Color) RepaColor {
	_, _, _, a := col.RGBA()
	cc, _ := colorful.MakeColor(col)

	return RepaColor{
		cc,
		float64(a) / 0xffff,
	}
}

func ParseColor(cstr string) (RepaColor, error) {
	if strings.HasPrefix(cstr, "#") {
		// hexa color
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
	} else if strings.HasPrefix(cstr, "rgb") {
		// rgb color
		var r, g, b uint8

		// transform commas/slashes to spaces to support legacy rgb() syntax
		cstr = strings.Replace(cstr, "rgba", "rgb", 1)
		cstr = strings.ReplaceAll(cstr, ",", " ")
		cstr = strings.Replace(cstr, "/", " ", 1)

		// CSS rgb() (simple)
		n, err := fmt.Sscanf(cstr, "rgb(%d %d %d)", &r, &g, &b)
		if err == nil && n == 3 {
			return RepaColor{
					colorful.Color{
						R: float64(r) / 255,
						G: float64(g) / 255,
						B: float64(b) / 255,
					},
					1.0,
				},
				nil
		}

		// CSS rgb() with alpha (as float)
		var a float64
		n, err = fmt.Sscanf(cstr, "rgb(%d %d %d %f)", &r, &g, &b, &a)
		if err == nil && n == 4 {
			return RepaColor{
					colorful.Color{
						R: float64(r) / 255,
						G: float64(g) / 255,
						B: float64(b) / 255,
					},
					a,
				},
				nil
		}

		// CSS rgb() with alpha (as percent)
		var aperc uint8
		n, err = fmt.Sscanf(cstr, "rgb(%d %d %d %d%%)", &r, &g, &b, &aperc)
		if err == nil && n == 4 {
			return RepaColor{
					colorful.Color{
						R: float64(r) / 255,
						G: float64(g) / 255,
						B: float64(b) / 255,
					},
					float64(aperc) / 100,
				},
				nil
		}

		// TODO `from` support
	}

	return NoColor, errors.New("cannot parse color")
}
