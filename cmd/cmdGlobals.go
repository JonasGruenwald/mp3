package cmd

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
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
