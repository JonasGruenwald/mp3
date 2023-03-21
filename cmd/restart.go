/*
Copyright © 2023 Jonas Grünwald
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "restart app_name|all",
	Short: "Restart a running service",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		forwardServiceCommand("restart", args)
	},
}

func init() {
	rootCmd.AddCommand(restartCmd)
}
