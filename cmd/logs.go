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
					runJournal(append([]string{"-u", serviceName}, args[1:]...))
				} else {
					// otherwise run default logs call on the service
					runJournal([]string{"-u", serviceName, "-n", "50", "-f"})
				}
			} else {
				// args are passed but are not service name - we pass the args on to journalctl for all services
				runJournal(append([]string{"-u", "mp3.*"}, args...))

			}
		} else {
			//runShell("journalctl", "-u", "mp3.*", "-n", "50", "-f", "-o", "short")
			runJournal([]string{"-u", "mp3.*", "-n", "50", "-f"})
		}

		if len(args) > 1 {
			runJournal(append([]string{"run", "-d"}, args...))
		} else if len(args) > 0 {
			runJournal([]string{"-u", getServiceName(args[0]), "-n", "50", "-f"})
		} else {
			runJournal([]string{"-u", "mp3.*", "-n", "50", "-f"})
		}
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
}
