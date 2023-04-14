/*
Copyright © 2023 Jonas Grünwald
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// reloadCmd represents the reload command
var reloadCmd = &cobra.Command{
	Use:   "reload service_name|all",
	Short: "Reload a running service (must be supported by app)",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		forwardServiceCommand("reload", args)
	},
}

func init() {
	rootCmd.AddCommand(reloadCmd)
}
