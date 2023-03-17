package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

const serviceNamePrefix = "pmu."

// const systemCtlUnitDir = "/etc/systemd/system"
const systemCtlUnitDir = "."

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

func runSilent(name string, args ...string) {
	cmd := exec.Command(name, args...)
	err := cmd.Run()
	if err != nil {
		fatal(fmt.Sprintf("Error running command %s %s", name, args))
	}
}

func runLoud(name string, args ...string) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	fmt.Printf("%s\n", output)
	if err != nil {
		fatal(fmt.Sprintf("Error running command %s %s", name, args))
	}
}
