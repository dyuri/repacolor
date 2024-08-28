package cmd

import (
	"bufio"
	"fmt"
	"log"

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
      fmt.Printf(`
  Distance:
    RGB:   %f
    Lab:   %f
    CIE76: %f
    CIE94: %f
    CIEDE: %f

`,
        refcolor.DistanceRgb(c.Color),
        refcolor.DistanceLab(c.Color),
        refcolor.DistanceCIE76(c.Color),
        refcolor.DistanceCIE94(c.Color),
        refcolor.DistanceCIEDE2000(c.Color),
      )

      // gradients
      // TODO use refcolor.BlendRgb(c.Color, .5) and friends
		}
	},
}

func init() {
	rootCmd.AddCommand(compareCmd)
}
