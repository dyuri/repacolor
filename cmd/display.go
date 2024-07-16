package cmd

import (
	"fmt"
	"image"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/mattn/go-isatty"

	"github.com/dyuri/repacolor/color"
)

var format string
var noansi bool

// TODO move to separate package ?
func renderAnsiImage(img image.Image) string {
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

		sb.WriteString("\n")
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

func getColorAnsiImage(c color.RepaColor, options ColorAnsiImageOptions) image.Image {
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

func textColorDetails(c color.RepaColor) string {
	nameStr, hasName := color.GetName(c)

	if hasName {
		nameStr = fmt.Sprintf("Name:  %s\n", nameStr)
	}

	return fmt.Sprintf("%sHex:   %s\nRGB:   %s\nHSL:   %s\nLAB:   %s\nLCH:   %s\nOKLAB: %s\nOKLCH: %s\n", nameStr, c.Hex(), c.RgbString(), c.HslString(), c.LabString(), c.LchString(), c.OkLabString(), c.OkLchString())
}

func mergeStringsVertically(a, b string) string {
	la := strings.Split(a, "\n")
	lb := strings.Split(b, "\n")
	// la = la[:len(la)-1]
	// lb = lb[:len(lb)-1]

	if len(la) < len(lb) {
		la = append(la, make([]string, len(lb)-len(la))...)
	} else if len(lb) < len(la) {
		lb = append(lb, make([]string, len(la)-len(lb))...)
	}

	var sb strings.Builder
	for i := 0; i < len(la); i++ {
		sb.WriteString(fmt.Sprintf("%s %s\n", la[i], lb[i]))
	}

	return sb.String()
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
		c, err := color.ParseColor(args[0], !nofallback)
		if err != nil {
			log.Fatal(err)
		}

		var repr string
		var termrepr string

		switch format {
		case "hex":
			repr = c.Hex()
		case "rgb", "rgba":
			repr = c.RgbString()
		case "hsl", "hsla":
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
		case "text":
			termrepr = textColorDetails(c)
		case "ansi":
			termrepr = renderAnsiImage(getColorAnsiImage(c, ColorAnsiImageOptions{}))
		default:
			ansirepr := renderAnsiImage(getColorAnsiImage(c, ColorAnsiImageOptions{}))
			textrepr := "\n" + textColorDetails(c)

			termrepr = mergeStringsVertically(ansirepr, textrepr)
		}

		// Print the color
		if isatty.IsTerminal(os.Stdout.Fd()) && !noansi {
			if termrepr == "" {
				termrepr = fmt.Sprintf("%s%s\033[0m\n", c.AnsiBg(), repr)
			}
			fmt.Print(termrepr)
		} else {
			if repr == "" {
				repr = c.Hex()
			}
			fmt.Printf("%s\n", repr)
		}
	},
}

func init() {
	displayCmd.Flags().StringVarP(&format, "format", "f", "", "Output format (hex, rgb, hsl, lab, lch, oklab, oklch)")
	displayCmd.Flags().BoolVarP(&noansi, "no-ansi", "n", false, "Disable ANSI color codes")

	rootCmd.AddCommand(displayCmd)
}
