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
	Short: "Show logs for services managed by pmu",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			runLoud("journalctl", append([]string{"run", "-d"}, args...)...)
		} else if len(args) > 0 {
			runLoud("journalctl", "-u", getServiceName(args[0]), "-n", "100", "-f", "-o", "short")
		} else {
			runLoud("journalctl", "-u", "pmu.*", "-n", "100", "-f", "-o", "short")
		}
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
}
