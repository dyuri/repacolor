package color

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
)

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
