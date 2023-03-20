/*
Copyright © 2023 Jonas Grünwald
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
	"github.com/struCoder/pidusage"
	"os"
	"strconv"
	"time"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "List all services managed by mp3",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		conn, err := dbus.NewSystemdConnectionContext(ctx)
		handleErr(err)
		units, err := conn.ListUnitsByPatternsContext(ctx, []string{}, []string{"mp3.*"})
		handleErr(err)

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(
			table.Row{
				text.FgCyan.Sprint("PID"),
				text.FgCyan.Sprint("Name"),
				text.FgCyan.Sprint("State"),
				text.FgCyan.Sprint("Uptime"),
				text.FgCyan.Sprint("Memory"),
				text.FgCyan.Sprint("CPU"),
			})
		for _, unit := range units {
			props, err := conn.GetAllPropertiesContext(ctx, unit.Name)
			handleErr(err)
			var memoryDisplay = ""
			var cpuDisplay = ""
			var uptimeDisplay = ""

			if unit.SubState == "running" {
				memoryCount := props["MemoryCurrent"].(uint64)
				memoryDisplay = fmt.Sprintf("%v", ByteCountSI(memoryCount))
				sysInfo, err := pidusage.GetStat(int(props["MainPID"].(uint32)))
				handleErr(err)
				cpuDisplay = fmt.Sprintf("%.2f%%", sysInfo.CPU)
				startTimeStamp, err := strconv.ParseInt(fmt.Sprintf("%v", props["ExecMainStartTimestamp"]), 10, 64)
				handleErr(err)
				fmt.Println(startTimeStamp)
				startTime := time.UnixMicro(startTimeStamp)
				uptimeDisplay = fmt.Sprintf("%v", time.Now().Sub(startTime).Truncate(time.Second))
			}

			t.AppendRow(table.Row{
				props["MainPID"],
				getAppName(unit.Name),
				colorStatus(unit.SubState),
				uptimeDisplay,
				memoryDisplay,
				cpuDisplay,
			})
		}

		t.SetStyle(table.StyleRounded)
		t.Render()
		conn.Close()
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
