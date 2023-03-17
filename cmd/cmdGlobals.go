package cmd

import (
	"errors"
	"fmt"
	"os"
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
