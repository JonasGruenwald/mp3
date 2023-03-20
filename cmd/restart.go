/*
Copyright © 2023 Jonas Grünwald
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart a running service",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] == "all" {
			fmt.Println("Restarting all services!")
			runShell("systemctl", "restart", "mp3.*")
		} else {
			var serviceName = getServiceName(args[0])
			runShell("systemctl", "restart", serviceName)
		}
	},
}

func init() {
	rootCmd.AddCommand(restartCmd)
}
