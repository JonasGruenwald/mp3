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

By default it will show the 50 most recent entries and tail the logs, 
however you can pass any arguments accepted by the journalctl command with --

Example:
# ge
mp3 logs -- -n 15
mp3 logs my_app -- --since 09:00 --until "1 hour ago"

`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			// check if the first arg is a service name
			if serviceExists(getServiceName(args[0])) {
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
			runJournal([]string{"-u", "mp3.*", "-n", "50", "-f"})
		}
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
}
