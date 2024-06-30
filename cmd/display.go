package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/dyuri/repacolor/color"
)

// displayCmd represents the display command
var displayCmd = &cobra.Command{
	Use:   "display",
	Args: cobra.ExactArgs(1),
	Short: "Display a color in the terminal",
	Long: `Display a color in the terminal.

Supported formats:
- Hexadecimal: #RRGGBB
- RGB: rgb(R G B) or rgb(R, G, B)`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := color.ParseColor(args[0])
		if err != nil {
			log.Fatal(err)
		}

		r := uint8(c.R * 255.0)
		g := uint8(c.G * 255.0)
		b := uint8(c.B * 255.0)

		// Print the color
		fmt.Printf("Color: \033[38;2;%d;%d;%dm%s\033[0m\n", r, g, b, c.Hex())
	},
}

func init() {
	rootCmd.AddCommand(displayCmd)
}
