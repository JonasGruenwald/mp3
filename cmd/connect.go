package cmd

import (
	"context"
	"fmt"
	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
	"os"
	"path"
	"strings"
	"text/template"
)

type CaddyFileConfig struct {
	Domain      string
	Port        string
	StaticPath  string
	RedirectWWW bool
}

func buildStaticCaddyFile(templateName string, config CaddyFileConfig, configFilePath string) {
	// get info & file modes
	caddyDirStat, err := os.Stat("/etc/caddy")
	handleErr(err)
	dirFileMode := caddyDirStat.Mode()

	// ensure dir
	err = os.MkdirAll(config.StaticPath, dirFileMode)
	fmt.Println("Created directory for static web assets:")
	fmt.Println(text.FgCyan.Sprint(config.StaticPath))
	handleErr(err)
	// construct caddyfile
	tmpl, err := template.ParseFS(TemplateFs, templateName)
	handleErr(err)
	file, err := os.Create(configFilePath)
	handleErr(err)
	err = tmpl.Execute(file, config)
	handleErr(err)
}

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect app_name|configuration domain_name",
	Short: "A brief description of your command",
	Long: `Connect a running mp3 app to your domain via a caddy reverse proxy.

Example:

mp3 connect my_server_app app.example.com

Instead of an app, you can also set up the following basic configurations:

SPA: Set up an SPA configuration for the specified domain
STATIC: Set up a static file server for the specified domain

Example:

mp3 connect STATIC example.com
`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		var config = CaddyFileConfig{
			RedirectWWW: false,
		}
		config.Domain = args[1]
		config.StaticPath = path.Join("/srv/www", config.Domain)
		if strings.Count(config.Domain, ".") == 0 {
			fatal("The domain argument must contain at least one '.'!")
		} else if strings.Count(config.Domain, ".") == 1 {
			config.RedirectWWW = true
		}
		var configFilePath = path.Join("/etc/caddy/sites", config.Domain+".conf")

		switch args[0] {
		case "SPA":
			{
				buildStaticCaddyFile("templates/spa-caddyfile.tmpl", config, configFilePath)
			}
		case "STATIC":
			{
				buildStaticCaddyFile("templates/static-caddyfile.tmpl", config, configFilePath)
			}
		default:
			{

				// Get App port
				ctx := context.Background()
				conn, err := dbus.NewSystemdConnectionContext(ctx)
				handleErrConn(err, conn)

				serviceName := getServiceName(args[0])
				// Ensure service exists
				if !serviceExists(serviceName) {
					fatal(fmt.Sprintf("Could not find service %s", serviceName))
				}
				props, err := conn.GetAllPropertiesContext(ctx, serviceName)
				handleErrConn(err, conn)
				conn.Close()
				ctx.Done()
				portMap := buildPortMap()
				appPorts := portMap[int(props["MainPID"].(uint32))]

				if strings.Count(appPorts, ",") > 0 {
					fmt.Println(fmt.Sprintf(
						"The service %s exposes more than one port, which port do you wish to connect?",
						text.Bold.Sprint(serviceName)))
					splitPorts := strings.Split(appPorts, ", ")
					config.Port = promptSelection(splitPorts)
				} else {
					config.Port = appPorts
				}
				fmt.Println(fmt.Sprintf("Connecting port %s of service %s to domain %s",
					text.Bold.Sprint(config.Port),
					text.Bold.Sprint(serviceName),
					text.Bold.Sprint(config.Domain),
				))
				tmpl, err := template.ParseFS(TemplateFs, "templates/reverse-proxy-caddyfile.tmpl")
				handleErr(err)
				file, err := os.Create(configFilePath)
				handleErr(err)
				err = tmpl.Execute(file, config)
				handleErr(err)
			}
		}

		fmt.Println("Created caddyfile config:")
		fmt.Println(text.FgCyan.Sprint(configFilePath))
		runSilent("systemctl", "reload", "caddy")
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
