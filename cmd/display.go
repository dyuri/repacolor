package cmd

import (
	"fmt"
	"image"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/mattn/go-isatty"

	"github.com/dyuri/repacolor/color"
)

var format string
var noansi bool

// TODO refactor
func displayAnsiImage(img image.Image) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	for y := 0; y < height; y += 2 {
		for x := 0; x < width; x++ {
			pixel := img.At(x, y)
			r, g, b, _ := pixel.RGBA()

			if y+1 < height {
				pixel2 := img.At(x, y+1)
				r2, g2, b2, _ := pixel2.RGBA()
				fmt.Printf("\033[48;2;%d;%d;%d;38;2;%d;%d;%d;1m▀\033[0m", r2>>8, g2>>8, b2>>8, r>>8, g>>8, b>>8)
			} else {
				fmt.Printf("\033[38;2;%d;%d;%d;1m▀\033[0m", r>>8, g>>8, b>>8)
			}

			if x == width-1 {
				fmt.Println()
			}
		}
	}
}

// displayCmd represents the display command
var displayCmd = &cobra.Command{
	Use:   "display <color>",
	Args: cobra.ExactArgs(1),
	Short: "Display a color in the terminal",
	Long: `Display a color in the terminal.

Supported formats:
- Hexadecimal: #RRGGBB
- RGB: rgb(R G B [/ A]) or rgb(R, G, B, A)`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := color.ParseColor(args[0])
		if err != nil {
			log.Fatal(err)
		}

		r := uint8(c.R * 255.0)
		g := uint8(c.G * 255.0)
		b := uint8(c.B * 255.0)

		var repr string

		switch format {
		case "hex":
			repr = c.Hex()
		case "rgb":
			repr = c.RgbString()
		case "hsl":
			repr = c.HslString()
		case "lab":
			repr = c.LabString()
		case "lch":
			repr = c.LchString()
		case "oklab":
			repr = c.OkLabString()
		case "oklch":
			repr = c.OkLchString()
		case "xyz":
			repr = c.XyzString()
		default:
			// TODO refactor

			img := image.NewRGBA(image.Rect(0, 0, 16, 16))
			for i := 0; i < 16; i++ {
				for j := 0; j < 16; j++ {
					pixel := color.LIGHTGRAY
					if (i+j)%2 != 0 {
						pixel = color.DARKGRAY
					}
					if i > 1 && i < 14 && j > 1 && j < 8 {
						pixel = c.AlphaBlendRgb(pixel, 2.2)
					}
					if i > 1 && i < 8 && j > 7 && j < 14 {
						pixel = c.AlphaBlendRgb(pixel, 1.0)
					}
					if i > 7 && i < 14 && j > 7 && j < 14 {
						co := c
						co.A = 1.0
						pixel = co.AlphaBlendRgb(pixel, 2.2)
					}
					img.Set(i, j, pixel)
				}
			}
			displayAnsiImage(img)

			repr = c.Hex()
		}

		// Print the color
		if isatty.IsTerminal(os.Stdout.Fd()) && !noansi {
			fgc := c.A11YPair()

			fr := uint8(fgc.R * 255.0)
			fg := uint8(fgc.G * 255.0)
			fb := uint8(fgc.B * 255.0)

			fmt.Printf("\033[48;2;%d;%d;%d;38;2;%d;%d;%d;1m%s\033[0m\n", r, g, b, fr, fg, fb, repr)
		} else {
			fmt.Printf("%s\n", repr)
		}
	},
}

func init() {
	displayCmd.Flags().StringVarP(&format, "format", "f", "", "Output format (hex, rgb, hsl, lab, lch, oklab, oklch)")
	displayCmd.Flags().BoolVarP(&noansi, "no-ansi", "n", false, "Disable ANSI color codes")

	rootCmd.AddCommand(displayCmd)
}
