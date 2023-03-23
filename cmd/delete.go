/*
Copyright © 2023 Jonas Grünwald
*/
package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a service definition that was created with mp3",
	Run: func(cmd *cobra.Command, args []string) {
		var serviceName = getServiceName(args[0])
		var targetServicePath = getServicePath(serviceName)

		runShell("systemctl", "stop", serviceName)
		runShell("systemctl", "disable", serviceName)
		runShell("systemctl", "reset-failed", serviceName)

		e := os.Remove(targetServicePath)
		if e != nil {
			fatal(e.Error())
		}

		runShell("systemctl", "daemon-reload")

	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
