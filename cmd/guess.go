package cmd

import (
	"github.com/spf13/cobra"
	"github.com/dyuri/repacolor/guess"
)

// pickCmd represents the pick command
var guessCmd = &cobra.Command{
	Use:   "guess",
	Args:  cobra.MaximumNArgs(1),
	Short: "Color guess game",
	Long: `Color guess game in the terminal.`,
	Run: func(cmd *cobra.Command, args []string) {
		guess.RunGuess()
	},
}

func init() {
	// TODO add flags for difficulity/etc.
	// guessCmd.Flags().BoolVarP(&showAlpha, "alpha", "a", false, "Show alpha channel")

	rootCmd.AddCommand(guessCmd)
}
