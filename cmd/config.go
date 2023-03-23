package cmd

import (
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config app_name",
	Short: "Edit the configuration for an MP3 managed app",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var serviceName = getServiceName(args[0])
		var targetServicePath = getServicePath(serviceName)
		if fileExists(targetServicePath) {
			runShell("nano", targetServicePath)
			runShell("systemctl", "daemon-reload")
		} else {
			fatal("Could not find service " + serviceName)
		}

	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
