package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bastjan/netstat"
	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func handleErr(err error) {
	if err != nil {
		fatal(err.Error())
	}
}

func handleErrConn(err error, conn *dbus.Conn) {
	if err != nil {
		conn.Close()
		fatal(err.Error())
	}
}

func fileExists(filePath string) bool {
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func printErrLn(message string) {
	fmt.Println(text.FgHiRed.Sprint(message))
}

func fatal(message string) {
	printErrLn(message)
	os.Exit(1)
}

func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func saveConfig() {
	homeDir, err := os.UserHomeDir()
	handleErr(err)
	err = viper.WriteConfigAs(path.Join(homeDir, "mp3-config.yaml"))
	handleErr(err)
}

func isAdoptedService(name string) bool {
	adoptedServices := viper.GetStringSlice("AdoptedServices")
	for _, service := range adoptedServices {
		if name+".service" == service {
			return true
		}
	}
	return false
}

func getServicePattern() []string {
	adoptedServices := viper.GetStringSlice("AdoptedServices")
	return append([]string{"mp3.*"}, adoptedServices...)
}

func getServiceName(target string) string {
	if isAdoptedService(target) {
		return target + ".service"
	}
	return fmt.Sprintf("%s%s.service", serviceNamePrefix, target)
}

// This function is just for constructing a path, not checking if the service path actually exists
func getServicePath(serviceName string) string {
	return path.Join(systemCtlUnitDir, serviceName)
}

func findServicePath(serviceName string) string {
	ctx := context.Background()
	conn, err := dbus.NewSystemdConnectionContext(ctx)
	handleErrConn(err, conn)
	unitFiles, err := conn.ListUnitFilesByPatternsContext(ctx, []string{}, []string{serviceName})
	handleErr(err)
	if len(unitFiles) > 1 {
		conn.Close()
		ctx.Done()
		fatal(fmt.Sprintf("More than one unit file found for service '%s', not sure what to do.", serviceName))
	} else if len(unitFiles) == 0 {
		conn.Close()
		ctx.Done()
		fatal(fmt.Sprintf("No unit file found for service '%s', not sure what to do.", serviceName))
	}
	conn.Close()
	ctx.Done()
	return unitFiles[0].Path
}

func getAppName(serviceName string) string {
	return strings.TrimSuffix(strings.TrimPrefix(serviceName, "mp3."), ".service")
}

func buildPortMap() map[int]string {
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
	return processPorts
}

func colorStatus(status string) string {
	switch status {
	case "running":
		return text.FgGreen.Sprint(status)
	case "dead":
		return text.FgRed.Sprint("stopped")
	case "failed":
		return text.FgRed.Sprint("failed")

	}
	return status
}

func colorEnabled(status string) string {
	if status == "enabled" {
		return text.Bold.Sprint(status)
	}
	return status
}

func humanByteCount(b uint64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func humanDuration(d time.Duration) string {
	if d < day {
		return d.String()
	}

	var b strings.Builder

	if d >= year {
		years := d / year
		_, err := fmt.Fprintf(&b, "%dy", years)
		handleErr(err)
		d -= years * year
	}

	days := d / day
	d -= days * day
	_, err := fmt.Fprintf(&b, "%dd%s", days, d)
	handleErr(err)

	return b.String()
}

func serviceExists(target string) bool {
	ctx := context.Background()
	conn, err := dbus.NewSystemdConnectionContext(ctx)
	handleErrConn(err, conn)
	unitFiles, err := conn.ListUnitFilesByPatternsContext(ctx, []string{}, []string{target})
	handleErr(err)
	conn.Close()
	ctx.Done()
	return len(unitFiles) > 0
}

func getConfirmation(question string) bool {
	if question == "" {
		question = "continue? [Y/N]"
	} else {
		question = question + " [Y/N]"
	}
	fmt.Println(question)
	reader := bufio.NewReader(os.Stdin)
	for {
		answer, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Read string failed, err: %v\n", err)
			break
		}
		answer = strings.ToLower(strings.TrimSpace(answer))
		if answer == "y" || answer == "yes" {
			return true
		} else if answer == "n" || answer == "no" {
			break
		} else {
			continue
		}
	}
	return false
}

func promptSelection(options []string) string {
	lower := 0
	upper := len(options) - 1
	for i, option := range options {
		fmt.Println(fmt.Sprintf("%v %s",
			text.FgCyan.Sprintf("[%v]", i),
			option,
		))
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		answer, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Read string failed, err: %v\n", err)
			break
		}
		answer = strings.TrimSpace(answer)
		if idx, err := strconv.Atoi(answer); err == nil && idx >= lower && idx <= upper {
			return options[idx]
		} else {
			fmt.Println(text.FgRed.Sprintf("Input must be between %v and %v",
				text.Bold.Sprint(lower),
				text.Bold.Sprint(upper)))
		}
	}
	return options[0]
}

func runSilent(name string, args ...string) {
	cmd := exec.Command(name, args...)
	err := cmd.Run()
	if err != nil {
		fatal(fmt.Sprintf("Error running command %s %s", name, args))
	}
}

func runSilentWithErr(name string, args ...string) {
	cmd := exec.Command(name, args...)
	err := cmd.Run()
	if err != nil {
		// ignore error, we press on
	}
}

func runShell(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fatal(fmt.Sprintf("Error running shell command %s %s", name, args))
	}
}

func getOutput(command string, args ...string) string {
	out, err := exec.Command(command, args...).Output()
	if err != nil {
		fatal(fmt.Sprintf("Error getting output from command %s %s", command, args))
	}
	cleaned := strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, string(out))
	return cleaned
}

func runJournal(args []string) {
	journalArgs := append(args, "-o", "json")
	cmd := exec.Command("journalctl", journalArgs...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	stdout, err := cmd.StdoutPipe()
	if nil != err {
		log.Fatalf(
			"Error obtaining stdout: %s", err.Error())
	}
	reader := bufio.NewReader(stdout)
	go func(reader io.Reader) {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			// TODO color text based on priority
			var journalEntry JournalEntry
			var err = json.Unmarshal([]byte(scanner.Text()), &journalEntry)
			handleErr(err)
			message := journalEntry.MESSAGE
			if journalEntry.PRIORITY != "6" {
				switch journalEntry.PRIORITY {
				case "0":
					fallthrough
				case "1":
					fallthrough
				case "2":
					message = text.FgRed.Sprint(journalEntry.MESSAGE)
				case "3":
					message = text.FgHiRed.Sprint(journalEntry.MESSAGE)
				case "4":
					message = text.FgHiYellow.Sprint(journalEntry.MESSAGE)
				case "5":
					message = text.Bold.Sprint(journalEntry.MESSAGE)
				case "7":
					message = text.Faint.Sprint(journalEntry.MESSAGE)
				}
			}

			if message == nil {
				message = text.FgYellow.Sprint("<message too long, pass -- --all to display anyways>")
			}
			timeStampUnix, err := strconv.ParseInt(journalEntry.Timestamp, 10, 64)
			handleErr(err)
			timeStamp := time.UnixMicro(timeStampUnix)
			title := ""
			if journalEntry.UNIT != "" {
				title = text.FgYellow.Sprint(getAppName(journalEntry.UNIT))
			} else if journalEntry.SystemdUnit != "" {
				title = text.FgCyan.Sprint(getAppName(journalEntry.SystemdUnit))
			}
			fmt.Println(fmt.Sprintf(
				"%s%s: %s",
				text.Faint.Sprint(timeStamp.Format("02.01.06 15:04")+"│"),
				title,
				message,
			))
		}
	}(reader)
	if err := cmd.Start(); nil != err {

		log.Fatalf("Error starting program: %s, %s", cmd.Path, err.Error())
	}
	cmd.Wait()
}

func forwardServiceCommand(cmd string, args []string) {
	if args[0] == "all" {
		runShell("systemctl", cmd, "mp3.*")
	} else {
		var serviceName = getServiceName(args[0])
		runShell("systemctl", cmd, serviceName)
	}
}
