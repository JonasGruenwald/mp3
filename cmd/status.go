/*
Copyright © 2023 Jonas Grünwald
*/
package cmd

import (
	"fmt"
	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/spf13/cobra"
	"path"
)

func unitFileLoaded(loadedUnits []dbus.UnitStatus, unitFile dbus.UnitFile) bool {
	for _, lu := range loadedUnits {
		if path.Join(systemCtlUnitDir, lu.Name) == unitFile.Path {
			return true
		}
		if path.Join(systemCtlUnitDirLib, lu.Name) == unitFile.Path {
			return true
		}
	}
	return false
}

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "List all services managed by mp3",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if !printStatus() {
			fmt.Println("No units to display, start an app with mp3 start!")
		}
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
