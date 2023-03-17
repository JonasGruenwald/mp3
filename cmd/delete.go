/*
Copyright © 2023 Jonas Grünwald
*/
package cmd

import (
	"os"
	"path"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a service definition that was created with pmu",
	Run: func(cmd *cobra.Command, args []string) {
		var serviceName = getServiceName(args[0])
		var targetServicePath = path.Join(systemCtlUnitDir, serviceName)

		runLoud("systemctl", "stop", serviceName)
		runLoud("systemctl", "disable", serviceName)

		e := os.Remove(targetServicePath)
		if e != nil {
			fatal(e.Error())
		}

		runLoud("systemctl", "daemon-reload")
		runLoud("systemctl", "reset-failed", serviceName)

	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
