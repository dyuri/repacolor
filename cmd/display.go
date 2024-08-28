package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
  "strings"

	"github.com/spf13/cobra"
	"github.com/mattn/go-isatty"

	"github.com/dyuri/repacolor/color"
	"github.com/dyuri/repacolor/display"
)

var format string
var noansi bool

// displayCmd represents the display command
var displayCmd = &cobra.Command{
	Use:   "display <color>",
	Short: "Display a color in the terminal",
	Long: `Display a color in the terminal.

Supported input formats:
- Hex: #RRGGBB
- RGB: rgb(R G B [/ A])
- HSL: hsl(H S% L% [/ A])

Supported output formats:
- Hex
- RGB
- HSL
- LAB
- LCH
- OKLAB
- OKLCH
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			// read from stdin
			inputReader := cmd.InOrStdin()
			scanner := bufio.NewScanner(inputReader)
			for scanner.Scan() {
				line := scanner.Text()
				args = append(args, line)
			}
		}
		for _, arg := range args {
			c, err := color.ParseColor(arg, !nofallback)
			if err != nil {
				log.Println(err)
				continue
			}

			var repr string
			var termrepr string

			switch strings.ToLower(format) {
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
				termrepr = display.TextColorDetails(c)
			case "ansi":
				termrepr = display.RenderAnsiImage(display.GetColorAnsiImage(c, display.ColorAnsiImageOptions{}))
			default:
				ansirepr := display.RenderAnsiImage(display.GetColorAnsiImage(c, display.ColorAnsiImageOptions{}))
				textrepr := "\n" + display.TextColorDetails(c)

				termrepr = display.MergeStringsVertically(ansirepr, textrepr, 24)
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
		}
	},
}

func init() {
	displayCmd.Flags().StringVarP(&format, "format", "f", "", "Output format (hex, rgb, hsl, lab, lch, oklab, oklch)")
	displayCmd.Flags().BoolVarP(&noansi, "no-ansi", "n", false, "Disable ANSI color codes")

	rootCmd.AddCommand(displayCmd)
}
