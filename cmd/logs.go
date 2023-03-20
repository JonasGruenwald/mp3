/*
Copyright © 2023 Jonas Grünwald
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs [service]",
	Short: "Show logs for services managed by mp3",
	Long: `Examples:
mp3 logs
mp3 logs <app_name>

You can pass on args directly to journalctl with --
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			// check if the first arg is a service name
			if serviceExists(args[0]) {
				var serviceName = getServiceName(args[0])
				// if any other args are passed, pass them on to journalctl
				if len(args) > 1 {
					runShell("journalctl", append([]string{"-u", serviceName}, args...)...)
				} else {
					// otherwise run default logs call on the service
					runShell("journalctl", "-u", serviceName, "-n", "50", "-f", "-o", "short")
				}
			} else {
				// args are passed but are not service name - we pass the args on to journalctl for all services
				runShell("journalctl", append([]string{"-u", "mp3.*"}, args...)...)

			}
		} else {
			runShell("journalctl", "-u", "mp3.*", "-n", "50", "-f", "-o", "short")
		}

		if len(args) > 1 {
			runShell("journalctl", append([]string{"run", "-d"}, args...)...)
		} else if len(args) > 0 {
			runShell("journalctl", "-u", getServiceName(args[0]), "-n", "50", "-f", "-o", "short")
		} else {
			runShell("journalctl", "-u", "mp3.*", "-n", "50", "-f", "-o", "short")
		}
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
}
