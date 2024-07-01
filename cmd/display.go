package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/dyuri/repacolor/color"
)

var format string

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

		switch format {
		case "hex":
			fmt.Println(c.Hex())
			return
		case "rgb":
			fmt.Println(c.RgbString())
			return
		}

		// Print the color
		fmt.Printf("Color: \033[38;2;%d;%d;%dm%s\033[0m (%s)\n", r, g, b, c.Hex(), format)
	},
}

func init() {
	displayCmd.Flags().StringVarP(&format, "format", "f", "", "Output format (hex, rgb)")

	rootCmd.AddCommand(displayCmd)
}
