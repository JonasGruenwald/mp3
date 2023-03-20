package cmd

import (
	"embed"
	"errors"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/text"
	"os"
	"os/exec"
	"path"
	"strings"
)

const serviceNamePrefix = "mp3."
const systemCtlUnitDir = "/etc/systemd/system"

var TemplateFs embed.FS

func handleErr(err error) {
	if err != nil {
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

func getAppName(serviceName string) string {
	return strings.TrimSuffix(strings.TrimPrefix(serviceName, "mp3."), ".service")
}

func colorStatus(status string) string {
	switch status {
	case "running":
		return text.FgGreen.Sprint(status)
	case "dead":
		return text.FgRed.Sprint("stopped")

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
	var targetServicePath = path.Join(systemCtlUnitDir, serviceName)
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
	err := cmd.Run() // add error checking
	if err != nil {
		fatal(fmt.Sprintf("Error running shell command %s %s", name, args))
	}
}
