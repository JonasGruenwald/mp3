package cmd

import (
	"context"
	"fmt"
	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/struCoder/pidusage"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Attempts to print the status but will return false if there are no units to display
func printStatus() bool {
	ctx := context.Background()
	conn, err := dbus.NewSystemdConnectionContext(ctx)
	handleErrConn(err, conn)
	units, err := conn.ListUnitsByPatternsContext(ctx, []string{}, getServicePattern())
	handleErrConn(err, conn)
	unitFiles, err := conn.ListUnitFilesByPatternsContext(ctx, []string{}, getServicePattern())
	handleErrConn(err, conn)

	// if there are no units to display, return false so a message can be displayed instead
	if len(units) == 0 && len(unitFiles) == 0 {
		return false
	}

	processPorts := buildPortMap()
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
			text.FgCyan.Sprint("User"),
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
		var userDisplay = ""

		if unit.SubState == "running" {
			pidDisplay = fmt.Sprintf("%v", props["MainPID"])
			processPort := processPorts[int(props["MainPID"].(uint32))]
			if processPort != "" {
				portDisplay = processPort
			}
			procStatus, err := os.ReadFile(fmt.Sprintf("/proc/%v/status", props["MainPID"]))
			if err != nil {
				fmt.Println(err)
			} else {
				re := regexp.MustCompile(`Uid:\s*(\d+)`)
				match := re.FindSubmatch(procStatus)
				if len(match) > 0 {
					uid := match[1]
					user, err := user.LookupId(string(uid))
					if err != nil {
						fmt.Println(err)
						userDisplay = text.FgYellow.Sprint("<?>")
					} else {
						userDisplay = user.Username
					}
				}
			}
			memoryCount := props["MemoryCurrent"].(uint64)
			memoryDisplay = fmt.Sprintf("%v", humanByteCount(memoryCount))
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
			userDisplay,
			colorEnabled(props["UnitFileState"].(string)),
		})
	}

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
			"",
			"",
			colorEnabled("disabled"),
		})
	}
	t.SetStyle(table.StyleRounded)
	t.Style().Options.SeparateRows = true
	t.Render()

	conn.Close()
	return true
}
