package color

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
)

// ParseColor //
////////////////

func TestParseHexColor(t *testing.T) {
	for i := 0; i < 100; i++ {
		r := rand.Uint32() % 256
		g := rand.Uint32() % 256
		b := rand.Uint32() % 256
		hexrep := fmt.Sprintf("#%02x%02x%02x", r, g, b)
		c, err := ParseColor(hexrep)
		if err != nil {
			t.Fatalf("Error parsing hexa color: %v", err)
		}
		if !strings.EqualFold(c.Hex(), hexrep) {
			t.Fatalf("Wrong color parsed: %v", c.Hex())
		}
		if c.A != 1 {
			t.Fatalf("Wrong alpha value: %v", c.A)
		}
		if !almosteq(c.R, float64(r)/0xff) || !almosteq(c.G, float64(g)/0xff) || !almosteq(c.B, float64(b)/0xff) {
			t.Fatalf("Wrong RGB values: %v %v %v (vs: %v %v %v)", c.R, c.G, c.B, r, g, b)
		}
	}
}

func TestParseHexaColor(t *testing.T) {
	for i := 0; i < 100; i++ {
		r := rand.Uint32() % 256
		g := rand.Uint32() % 256
		b := rand.Uint32() % 256
		a := rand.Uint32() % 256
		hexrep := fmt.Sprintf("#%02x%02x%02x%02x", r, g, b, a)
		c, err := ParseColor(hexrep)
		if err != nil {
			t.Fatalf("Error parsing hexa color: %v", err)
		}
		if a != 255 && !strings.EqualFold(c.Hex(), hexrep) {
			t.Fatalf("Wrong color parsed: %v", c.Hex())
		} else if a == 255 && !strings.EqualFold(c.Hex(), fmt.Sprintf("#%02x%02x%02x", r, g, b)) {
			t.Fatalf("Wrong color parsed: %v (alpha is 1)", c.Hex())
		}
		if !almosteq(c.R, float64(r)/0xff) || !almosteq(c.G, float64(g)/0xff) || !almosteq(c.B, float64(b)/0xff) || !almosteq(c.A, float64(a)/0xff) {
			t.Fatalf("Wrong RGB values: %v %v %v %v (vs: %v %v %v %v)", c.R, c.G, c.B, c.A, r, g, b, a)
		}
	}
}

func TestParseRgbColor(t *testing.T) {
	for i := 0; i < 100; i++ {
		r := rand.Uint32() % 256
		g := rand.Uint32() % 256
		b := rand.Uint32() % 256
		rgbrep := fmt.Sprintf("rgb(%d, %d, %d)", r, g, b)
		rgbrep2 := fmt.Sprintf("rgb(%d %d %d)", r, g, b)
		c, err := ParseColor(rgbrep)
		if err != nil {
			t.Fatalf("Error parsing rgb color: %v", err)
		}
		c2, err := ParseColor(rgbrep2)
		if err != nil {
			t.Fatalf("Error parsing rgb color: %v", err)
		}
		if !strings.EqualFold(c.Hex(), fmt.Sprintf("#%02x%02x%02x", r, g, b)) {
			t.Fatalf("Wrong color parsed: %v", c.Hex())
		}
		if c.A != 1 {
			t.Fatalf("Wrong alpha value: %v", c.A)
		}
		if !almosteq(c.R, float64(r)/0xff) || !almosteq(c.G, float64(g)/0xff) || !almosteq(c.B, float64(b)/0xff) {
			t.Fatalf("Wrong RGB values: %v %v %v (vs: %v %v %v)", c.R, c.G, c.B, r, g, b)
		}
		if c != c2 {
			t.Fatalf("Error parsing rgb variants: %v, %v", rgbrep, rgbrep2)
		}
	}
}

func TestParseRgbaColor(t *testing.T) {
	for i := 0; i < 100; i++ {
		r := rand.Uint32() % 256
		g := rand.Uint32() % 256
		b := rand.Uint32() % 256
		a := rand.Float64()
		rgbrep := fmt.Sprintf("rgba(%d, %d, %d, %f)", r, g, b, a)
		rgbrep2 := fmt.Sprintf("rgb(%d %d %d / %f)", r, g, b, a)
		rgbrep3 := fmt.Sprintf("rgb(%d %d %d / %d%%)", r, g, b, uint8(a * 100))
		c, err := ParseColor(rgbrep)
		if err != nil {
			t.Fatalf("Error parsing rgba color: %v", err)
		}
		c2, err := ParseColor(rgbrep2)
		if err != nil {
			t.Fatalf("Error parsing rgba color: %v", err)
		}
		c3, err := ParseColor(rgbrep2)
		if err != nil {
			t.Fatalf("Error parsing rgba color: %v", err)
		}
		if !almosteq(c.R, float64(r)/0xff) || !almosteq(c.G, float64(g)/0xff) || !almosteq(c.B, float64(b)/0xff) || !almosteq(c.A, a) {
			t.Fatalf("Wrong RGB values: %v %v %v %v (vs: %v %v %v %v)", c.R, c.G, c.B, c.A, r, g, b, a)
		}
		if c != c2 {
			t.Fatalf("Error parsing rgba variants: %v, %v", rgbrep, rgbrep2)
		}
		if c != c3 {
			t.Fatalf("Error parsing rgba variants: %v, %v", rgbrep, rgbrep3)
		}
	}
}

// switched to mazznoer/csscolorparser, so this test is no longer needed, but I keep it none the less
