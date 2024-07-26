/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/dyuri/repacolor/color"
	"github.com/dyuri/repacolor/picker"
)

var mode string
var useAlpha bool

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

		fmt.Printf("pick called, mode: %s, color: %s\n", mode, c.Hex())
		picker.RunPicker(c)
	},
}

func init() {
	displayCmd.Flags().StringVarP(&mode, "mode", "m", "rgb", "Output format (rgb, hsl, lab, lch, oklab, oklch)")
	displayCmd.Flags().BoolVarP(&useAlpha, "alpha", "a", false, "Use alpha channel")

	rootCmd.AddCommand(pickCmd)
}
