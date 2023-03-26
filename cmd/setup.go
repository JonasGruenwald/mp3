/*
Copyright © 2023 Jonas Grünwald
*/
package cmd

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
	"os"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up specific apps to work with mp3",
	Long: `Available options:

setup caddy

Sets up the caddy webserver to let mp3 automatically connect apps to domains.
Caddy must be installed first, see:
https://caddyserver.com/docs/install
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "caddy":
			{
				fmt.Println(`MP3 will create a new caddyfile structure under /etc/caddy,
your existing /etc/caddy/Caddyfile will be moved.`)
				if getConfirmation("Are you sure you want to do this?") {
					caddyDirStat, err := os.Stat("/etc/caddy")
					handleErr(err)
					dirFileMode := caddyDirStat.Mode()
					caddyFileStat, err := os.Stat("/etc/caddy/Caddyfile")
					handleErr(err)
					fileFileMode := caddyFileStat.Mode()
					// structuring calls
					err = os.Mkdir("/etc/caddy/sites", dirFileMode)
					handleErr(err)
					err = os.Rename("/etc/caddy/Caddyfile", "/etc/caddy/sites/default.conf")
					handleErr(err)
					err = os.WriteFile("/etc/caddy/Caddyfile", []byte("import /etc/caddy/sites/*.conf"), fileFileMode)
					handleErr(err)
					// info printing
					fmt.Println("Structure, has been set up successfully, here is what your /etc/caddy structure looks like now:")
					fmt.Println("")
					l := list.NewWriter()
					l.SetOutputMirror(os.Stdout)
					l.AppendItem("/etc/caddy")
					l.Indent()
					l.AppendItem(fmt.Sprintf("Caddyfile   %s", text.FgYellow.Sprint("← New Caddyfile just imports sites")))
					l.AppendItem("/sites")
					l.Indent()
					l.AppendItems([]interface{}{
						fmt.Sprintf("default.conf      %s", text.FgYellow.Sprint("← Your old Caddyfile")),
						text.FgHiCyan.Sprint("site.b.com.conf"),
						fmt.Sprintf("%s   %s",
							text.FgHiCyan.Sprint("site.b.com.conf"),
							text.FgYellow.Sprint("← mp3 will create conf files for your apps like this"),
						),
						text.FgHiCyan.Sprint("site.c.com.conf"),
						fmt.Sprintf("%s    %s",
							text.FgHiCyan.Sprint("other.net.conf"),
							text.FgYellow.Sprint("← add extra conf files for static sites etc."),
						),
					})
					l.SetStyle(list.StyleConnectedRounded)
					l.Render()
					fmt.Println("")
					fmt.Println("Run 'mp3 connect app_name domain_name' to connect your apps.")
					fmt.Println("")
					runSilentWithErr("systemctl", "caddy", "reload")
				}
			}
		default:
			fatal(fmt.Sprintf("Don't know how to set up %s", args[0]))
		}
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
