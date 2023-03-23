package cmd

import (
	"embed"
	"time"
)

const serviceNamePrefix = "mp3."
const systemCtlUnitDir = "/etc/systemd/system"

const (
	day  = time.Minute * 60 * 24
	year = 365 * day
)

var TemplateFs embed.FS

type JournalEntry struct {
	MESSAGE          interface{}
	PRIORITY         string
	SyslogIdentifier string `json:"SYSLOG_IDENTIFIER"`
	UNIT             string
	SystemdUnit      string `json:"_SYSTEMD_UNIT"`
	CmdLine          string `json:"_CMDLINE"`
	Timestamp        string `json:"__REALTIME_TIMESTAMP"`
}
