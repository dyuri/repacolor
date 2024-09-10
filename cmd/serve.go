package cmd

import (
	"github.com/spf13/cobra"
	"github.com/dyuri/repacolor/picker"
)

var serveCmd = &cobra.Command{
	Use:   "serve [port]",
	Args:  cobra.MaximumNArgs(1),
	Short: "Serve picker through ssh",
	Long: `Interactive color picker for the terminal.`,
	Run: func(cmd *cobra.Command, args []string) {
		var sshPort string
		if len(args) == 0 {
			sshPort = "10022"
		} else {
			sshPort = args[0]
		}
		picker.ServePicker(sshPort)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
