/*
Copyright © 2023 Jonas Grünwald
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/bastjan/netstat"
	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
	"github.com/struCoder/pidusage"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func unitFileLoaded(loadedUnits []dbus.UnitStatus, unitFile dbus.UnitFile) bool {
	for _, lu := range loadedUnits {
		if path.Join(systemCtlUnitDir, lu.Name) == unitFile.Path {
			return true
		}
	}
	return false
}

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "List all services managed by mp3",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		conn, err := dbus.NewSystemdConnectionContext(ctx)
		handleErrConn(err, conn)
		units, err := conn.ListUnitsByPatternsContext(ctx, []string{}, []string{"mp3.*"})
		handleErrConn(err, conn)

		// Get port map
		processPorts := make(map[int]string)
		connections, err := netstat.TCP.Connections()
		handleErr(err)
		connections6, err := netstat.TCP6.Connections()
		handleErr(err)
		connections = append(connections, connections6...)
		for _, connection := range connections {
			if processPorts[connection.Pid] == "" {
				processPorts[connection.Pid] = fmt.Sprintf("%v", connection.Port)
			} else {
				if !strings.Contains(processPorts[connection.Pid], strconv.Itoa(connection.Port)) {
					processPorts[connection.Pid] = fmt.Sprintf("%s, %v",
						processPorts[connection.Pid],
						connection.Port)
				}
			}
		}

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
				text.FgCyan.Sprint("Ports"),
				text.FgCyan.Sprint("Startup"),
			})
		for _, unit := range units {
			props, err := conn.GetAllPropertiesContext(ctx, unit.Name)
			handleErrConn(err, conn)
			var pidDisplay = ""
			var memoryDisplay = ""
			var cpuDisplay = ""
			var uptimeDisplay = ""
			var portDisplay = ""

			if unit.SubState == "running" {
				pidDisplay = fmt.Sprintf("%v", props["MainPID"])
				processPort := processPorts[int(props["MainPID"].(uint32))]
				if processPort != "" {
					portDisplay = processPort
				}
				memoryCount := props["MemoryCurrent"].(uint64)
				memoryDisplay = fmt.Sprintf("%v", ByteCountSI(memoryCount))
				sysInfo, err := pidusage.GetStat(int(props["MainPID"].(uint32)))
				handleErrConn(err, conn)
				cpuDisplay = fmt.Sprintf("%.2f%%", sysInfo.CPU)
				startTimeStamp, err := strconv.ParseInt(fmt.Sprintf("%v", props["ExecMainStartTimestamp"]), 10, 64)
				handleErrConn(err, conn)
				startTime := time.UnixMicro(startTimeStamp)
				uptimeDisplay = humanDuration(time.Now().Sub(startTime).Truncate(time.Second))
			}

			t.AppendRow(table.Row{
				pidDisplay,
				getAppName(unit.Name),
				colorStatus(unit.SubState),
				uptimeDisplay,
				memoryDisplay,
				cpuDisplay,
				portDisplay,
				colorEnabled(props["UnitFileState"].(string)),
			})
		}

		unitFiles, err := conn.ListUnitFilesByPatternsContext(ctx, []string{}, []string{"mp3.*"})
		handleErrConn(err, conn)

		var deadUnits []string
		for _, unitFile := range unitFiles {
			if !unitFileLoaded(units, unitFile) {
				deadUnits = append(deadUnits,
					strings.TrimPrefix(strings.TrimSuffix(filepath.Base(unitFile.Path),
						filepath.Ext(unitFile.Path),
					), "mp3."))
			}
		}
		for _, deadUnit := range deadUnits {
			t.AppendRow(table.Row{
				"",
				deadUnit,
				colorStatus("dead"),
				"",
				"",
				"",
				colorEnabled("disabled"),
			})
		}
		t.SetStyle(table.StyleRounded)
		t.Style().Options.SeparateRows = true
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
