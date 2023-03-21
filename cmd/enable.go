/*
Copyright © 2023 Jonas Grünwald
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// enableCmd represents the enable command
var enableCmd = &cobra.Command{
	Use:   "enable app_name|all",
	Short: "Enable a service (run on startup)",
	Run: func(cmd *cobra.Command, args []string) {
		forwardServiceCommand("enable", args)
		runShell("systemctl", "daemon-reload")
	},
}

func init() {
	rootCmd.AddCommand(enableCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// enableCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// enableCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
