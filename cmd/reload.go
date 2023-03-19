/*
Copyright © 2023 Jonas Grünwald
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// reloadCmd represents the reload command
var reloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "Reload a running service (must be supported by app)",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] == "all" {
			runLoud("systemctl", "reload", "*.mp3")
		} else {
			var serviceName = getServiceName(args[0])
			runLoud("systemctl", "reload", serviceName)
		}
	},
}

func init() {
	rootCmd.AddCommand(reloadCmd)
}
