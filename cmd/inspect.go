/*
Copyright © 2023 Jonas Grünwald
*/
package cmd

import (
	"context"
	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
	"os"
)

// inspectCmd represents the inspect command
var inspectCmd = &cobra.Command{
	Use:   "inspect [app_name]",
	Short: "Dump systemd information about a given app running as a service (for debugging)",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var serviceName = getServiceName(args[0])
		ctx := context.Background()
		conn, err := dbus.NewSystemdConnectionContext(ctx)
		handleErr(err)
		props, err := conn.GetAllPropertiesContext(ctx, serviceName)
		handleErr(err)
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(
			table.Row{
				text.FgCyan.Sprint("Property"),
				text.FgCyan.Sprint("Value"),
			})
		for key, element := range props {
			t.AppendRow(table.Row{
				key,
				element,
			})
		}
		t.SetStyle(table.StyleLight)
		t.Style().Options.SeparateRows = true

		t.Render()
		conn.Close()
		conn.Close()
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// inspectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// inspectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
