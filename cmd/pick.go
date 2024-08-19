package cmd

import (
	"github.com/spf13/cobra"
	"github.com/dyuri/repacolor/color"
	"github.com/dyuri/repacolor/picker"
)

var showAlpha bool

// pickCmd represents the pick command
var pickCmd = &cobra.Command{
	Use:   "pick [color]",
	Args:  cobra.MaximumNArgs(1),
	Short: "Color picker",
	Long: `Interactive color picker for the terminal.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := color.WHITE
		if (len(args) > 0) {
			c, _ = color.ParseColor(args[0], true)
		}

		picker.RunPicker(c, showAlpha)
	},
}

func init() {
	pickCmd.Flags().BoolVarP(&showAlpha, "alpha", "a", false, "Show alpha channel")

	rootCmd.AddCommand(pickCmd)
}
