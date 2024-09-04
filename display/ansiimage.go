package display

import (
	"image"
	"fmt"
	"strings"

	"github.com/dyuri/repacolor/color"
)

func RenderAnsiImage(img image.Image) string {
	var sb strings.Builder
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	for y := 0; y < height; y += 2 {
		for x := 0; x < width; x++ {
			pixel := img.At(x, y)
			r, g, b, a := pixel.RGBA()

			if y+1 < height {
				pixel2 := img.At(x, y+1)
				r2, g2, b2, a2 := pixel2.RGBA()
				if a > 0 && a2 > 0 {
					sb.WriteString(fmt.Sprintf("\033[48;2;%d;%d;%d;38;2;%d;%d;%d;1m▀\033[0m", r2>>8, g2>>8, b2>>8, r>>8, g>>8, b>>8))
				} else if a > 0 {
					sb.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%d;1m▀\033[0m", r>>8, g>>8, b>>8))
				} else if a2 > 0 {
					sb.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%d;1m▄\033[0m", r2>>8, g2>>8, b2>>8))
				} else {
					sb.WriteString(" ");
				}
			} else {
				sb.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%d;1m▀\033[0m", r>>8, g>>8, b>>8))
			}

		}

    if y + 2 < height {
      sb.WriteString("\n")
    }
	}

	return sb.String()
}

type ColorAnsiImageOptions struct {
	Width int
	Height int
	Margin int
	Padding int
	Background1 color.RepaColor
	Background2 color.RepaColor
}

func GetColorAnsiImage(c color.RepaColor, options ColorAnsiImageOptions) image.Image {
	if (options.Width == 0) {
		options.Width = 16
	}
	if (options.Height == 0) {
		options.Height = 16
	}
	if (options.Margin == 0) {
		options.Margin = 2
	}
	if (options.Padding == 0) {
		options.Padding = 2
	}
	if (options.Background1 == color.NOCOLOR) {
		options.Background1 = color.LIGHTGRAY
	}
	if (options.Background2 == color.NOCOLOR) {
		options.Background2 = color.DARKGRAY
	}

	fullwidth := options.Width + 2 * options.Margin + 2 * options.Padding
	fullheight := options.Height + 2 * options.Margin + 2 * options.Padding
	mp := options.Margin + options.Padding

	img := image.NewRGBA(image.Rect(0, 0, fullwidth, fullheight))
	for j := 0; j < fullheight; j++ {
		for i := 0; i < fullwidth; i++ {
			pixel := color.NOCOLOR
			if i >= options.Margin && i < fullheight - options.Margin && j >= options.Margin && j < fullwidth - options.Margin {
				pixel = options.Background1
				if (i+j)%2 != 0 {
					pixel = options.Background2
				}
			}
			if i >= mp  && i < fullheight - mp && j >= mp && j < fullwidth / 2 {
				pixel = c.AlphaBlendRgb(pixel, 2.2)
			}
			if i >= mp && i < fullheight / 2 && j >= fullwidth / 2 && j < fullwidth - mp {
				pixel = c.AlphaBlendRgb(pixel, 1.0)
			}
			if i >= fullheight / 2 && i < fullheight - mp && j >= fullwidth / 2 && j < fullwidth - mp {
				co := c
				co.A = 1.0
				pixel = co.AlphaBlendRgb(pixel, 2.2)
			}
			img.Set(i, j, pixel)
		}
	}

	return img
}

func GetCompareAnsiImage(c1, c2 color.RepaColor, options ColorAnsiImageOptions) image.Image {
	if (options.Width == 0) {
		options.Width = 16
	}
	if (options.Height == 0) {
		options.Height = 16
	}
	if (options.Margin == 0) {
		options.Margin = 2
	}
	if (options.Padding == 0) {
		options.Padding = 2
	}
	if (options.Background1 == color.NOCOLOR) {
		options.Background1 = color.LIGHTGRAY
	}
	if (options.Background2 == color.NOCOLOR) {
		options.Background2 = color.DARKGRAY
	}

	fullwidth := options.Width + 2 * options.Margin + 2 * options.Padding
	fullheight := options.Height + 2 * options.Margin + 2 * options.Padding
	mp := options.Margin + options.Padding

	img := image.NewRGBA(image.Rect(0, 0, fullwidth, fullheight))
	for i := 0; i < fullheight; i++ {
		for j := 0; j < fullwidth; j++ {
			pixel := color.NOCOLOR
			if i >= options.Margin && i < fullheight - options.Margin && j >= options.Margin && j < fullwidth - options.Margin {
				pixel = options.Background1
				if (i+j)%2 != 0 {
					pixel = options.Background2
				}
			}

      c := c1
      if i >= mp && i >= fullwidth / 2 {
        c = c2
      }

      if j >= mp && j < fullheight - mp && i >= mp && i < fullwidth - mp {
        pixel = c.AlphaBlendRgb(pixel, 2.2)
        if j >= fullheight / 2 && j < fullheight - mp {
          co := c
          co.A = 1.0
          pixel = co.AlphaBlendRgb(pixel, 2.2)
        }
			}
			img.Set(i, j, pixel)
		}
	}

	return img
}

func AnsiGradient(c1, c2 color.RepaColor, width, mode int) string {
	var sb strings.Builder
	for i := 0; i < width; i++ {
		c := c1.Blend(c2, float64(i) / float64(width - 1), mode, false)
		sb.WriteString(c.AnsiBg())
		sb.WriteString(" ")
	}
	sb.WriteString(color.ANSI_RESET)
	return sb.String()
}
