package color

import (
	"testing"

	"fmt"
	"math/rand"
	"strconv"

	"github.com/lucasb-eyer/go-colorful"
)

func TestRgbStringNoAlpha(t *testing.T) {
	for i := 0; i < 100; i++ {
		r := rand.Uint32() % 256
		g := rand.Uint32() % 256
		b := rand.Uint32() % 256
		c := RepaColor{Color: colorful.Color{R: float64(r) / 0xff, G: float64(g) / 0xff, B: float64(b) / 0xff}, A: 1.0}

		s := c.RgbString()

		if s != fmt.Sprintf("rgb(%d %d %d)", r, g, b) {
			t.Fatalf("Invalid RGB representation: %v (vs. %d %d %d)", s, r, g, b)
		}
	}
}

func TestRgbStringAlpha(t *testing.T) {
	for i := 0; i < 100; i++ {
		r := rand.Uint32() % 256
		g := rand.Uint32() % 256
		b := rand.Uint32() % 256
		a := rand.Float64()
		c := RepaColor{Color: colorful.Color{R: float64(r) / 0xff, G: float64(g) / 0xff, B: float64(b) / 0xff}, A: a}

		s := c.RgbString()
		arep := strconv.FormatFloat(a, 'f', -1, 64)

		if s != fmt.Sprintf("rgb(%d %d %d / %s)", r, g, b, arep) {
			t.Fatalf("Invalid RGB representation: %v (vs. %d %d %d %s)", s, r, g, b, arep)
		}
	}
}
