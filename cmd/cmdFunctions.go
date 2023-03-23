package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/jedib0t/go-pretty/v6/text"
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

func fatal(message string) {
	fmt.Println(message)
	os.Exit(1)
}

func getServiceName(target string) string {
	return fmt.Sprintf("%s%s.service", serviceNamePrefix, target)
}

func getServicePath(serviceName string) string {
	return path.Join(systemCtlUnitDir, serviceName)
}

func getAppName(serviceName string) string {
	return strings.TrimSuffix(strings.TrimPrefix(serviceName, "mp3."), ".service")
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

func ByteCountSI(b uint64) string {
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

func serviceExists(target string) bool {
	var serviceName = getServiceName(target)
	var targetServicePath = getServicePath(serviceName)
	return fileExists(targetServicePath)
}

func runSilent(name string, args ...string) {
	cmd := exec.Command(name, args...)
	err := cmd.Run()
	if err != nil {
		fatal(fmt.Sprintf("Error running command %s %s", name, args))
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

func humanDuration(d time.Duration) string {
	if d < day {
		return d.String()
	}

	var b strings.Builder

	if d >= year {
		years := d / year
		fmt.Fprintf(&b, "%dy", years)
		d -= years * year
	}

	days := d / day
	d -= days * day
	fmt.Fprintf(&b, "%dd%s", days, d)

	return b.String()
}

func runJournal(args []string) {
	cmd := exec.Command("journalctl", append(args, "-o", "json")...)
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
				"%s|%s|%s:",
				text.Faint.Sprint(timeStamp.Format("02.01.06 15:04")),
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
