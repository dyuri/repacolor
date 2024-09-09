package cmd

import (
	"bufio"
	"fmt"
	"log"

	"golang.org/x/term"

	"github.com/spf13/cobra"

	"github.com/dyuri/repacolor/color"
	"github.com/dyuri/repacolor/display"
)

var compareCmd = &cobra.Command{
	Use:   "compare <color1> <color2>",
	Short: "Compare the given colors",
	Long: `Compare the given colors in the terminal.

For supported input formats, see the 'display' command.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			// read from stdin
			inputReader := cmd.InOrStdin()
			scanner := bufio.NewScanner(inputReader)
			for scanner.Scan() {
				line := scanner.Text()
				args = append(args, line)
			}
		}
		if len(args) < 2 {
			log.Fatal("At least two colors are required")
		}

		refcolor, err := color.ParseColor(args[0], !nofallback)
		if err != nil {
			log.Fatal(err)
		}

		// compare colors with the first one
		for _, arg := range args[1:] {
			c, err := color.ParseColor(arg, !nofallback)
			if err != nil {
				log.Println(err)
				continue
			}

			ansirepr := display.RenderAnsiImage(display.GetCompareAnsiImage(refcolor, c, display.ColorAnsiImageOptions{}))
			textrepr1 := "\n" + display.TextColorDetails(refcolor)
			textrepr2 := "\n" + display.TextColorDetails(c)
			textrepr := display.MergeStringsVertically(textrepr1, textrepr2, 0)

			termrepr := display.MergeStringsVertically(ansirepr, textrepr, 24)

			fmt.Print(termrepr)

			// distances
			fmt.Printf("\n  Distance:\n    RGB: %f CIE76: %f CIE94: %f CIEDE: %f\n\n",
				refcolor.DistanceRgb(c.Color),
				refcolor.DistanceCIE76(c.Color),
				refcolor.DistanceCIE94(c.Color),
				refcolor.DistanceCIEDE2000(c.Color),
			)

			// gradients
			terminalWidth, _, _ := term.GetSize(0)
			if terminalWidth <= 4 {
				terminalWidth = 80
			}

			for _, mode := range []int{
					color.BLEND_RGB,
					color.BLEND_LINEARRGB,
					color.BLEND_HSV,
					color.BLEND_LAB,
					color.BLEND_OKLAB,
					color.BLEND_LCH,
					color.BLEND_OKLCH,
					// color.BLEND_XYZ,
				} {
				grad := display.AnsiGradient(refcolor, c, terminalWidth - 4, mode)
				fmt.Printf("  %s\n", grad)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(compareCmd)
}
