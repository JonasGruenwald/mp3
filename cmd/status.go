/*
Copyright © 2023 Jonas Grünwald
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "List all services managed by mp3",
	Run: func(cmd *cobra.Command, args []string) {
		runLoud("systemctl", "list-units", "mp3.*", "--full", "--all", "--no-pager")
	},
}

var statusAliasList = &cobra.Command{
	Use:   "list",
	Short: "alias for status",
	Run: func(cmd *cobra.Command, args []string) {
		statusCmd.Run(cmd, args)
	},
}

var statusAliasLs = &cobra.Command{
	Use:   "ls",
	Short: "alias for status",
	Run: func(cmd *cobra.Command, args []string) {
		statusCmd.Run(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(statusAliasList)
	rootCmd.AddCommand(statusAliasLs)
}
