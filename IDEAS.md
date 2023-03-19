# Status

Could implement a prettier version of `status` by using `go-systemd/dbus` to parse running units

```go
package cmd

import (
	"github.com/coreos/go-systemd/dbus"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	conn, err, := dbus.NewSystemdConnection()
	handleErr(err)
	units, err := conn.ListUnits()
	handleErr(err)
	for i, unit := range units{
	// parse unit info
}
}

```

And then use something like https://github.com/jedib0t/go-pretty to display them in a nice table.

# Connect
Special Code for starting caddy service with its recommended service file, and setting up a caddyfile structure that 
allows dynamically adding sites by creating new files in a folder, files are automatically imported caddyfiles.
```shell
mp3 start caddy
```

magic command to connect a mp3 service to a domain, by detecting its port and creating a caddyfile config for it

```go
mp3 connect app app.example.com
```
