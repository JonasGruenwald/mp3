package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

type ServiceSettings struct {
	AppName           string
	RestartDelay      int
	NoAutorestart     bool
	Interpreter       string
	UserName          string
	WorkingDir        string
	ExecStart         string
	CreateServiceOnly bool
}

var settings = ServiceSettings{}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start entrypoint|app_name|all",
	Short: "Start an application with mp3",
	Long: `Supply the name of an application entry point in your current working directory,
and mp3 will create a daemonized service for it for you.
Pass the name of an existing mp3-created service to start it

Examples:

# Spin up new service:
mp3 start app.js
mp3 start bashscript.sh
mp3 start python-app.py
mp3 start ./binary-file -- --port 1520

# Start existing service
mp3 start my-app

`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		/**
		We need to figure out if the user wants to spin up a service, or start an existing service.
		We will check which option is possible and decide based on that
		*/
		target := args[0]
		if target == "all" {
			fmt.Println("Starting all services!")
			runShell("systemctl", "start", "mp3.*", "--all")
			return
		}
		workingDir, err := os.Getwd()
		if err != nil {
			fatal("Could not get working directory")
		}
		// path to potential new app to create service for
		var targetPath = target
		if !path.IsAbs(targetPath) {
			targetPath, err = filepath.Abs(target)
			handleErr(err)
		}
		// path to potential existing service to start
		if settings.AppName == "" {
			settings.AppName = strings.TrimSuffix(filepath.Base(targetPath), filepath.Ext(targetPath))
		}
		var serviceName = getServiceName(settings.AppName)
		if settings.CreateServiceOnly || !serviceExists(serviceName) {
			// Check first that the target file exists
			if fileExists(targetPath) {
				fatal("Can't find file: " + targetPath)
			}
			// We want to create a new service
			settings.ExecStart = targetPath
			// If we are dealing with a script file, we need to add the interpreter
			if strings.HasSuffix(target, ".js") {
				if settings.Interpreter == "" {
					nodePath, err := exec.LookPath("node")
					handleErr(err)
					settings.Interpreter = nodePath
				}
				settings.ExecStart = settings.Interpreter + " " + targetPath
			} else if strings.HasSuffix(target, ".py") {
				if settings.Interpreter == "" {
					pythonPath, err := exec.LookPath("node")
					handleErr(err)
					settings.Interpreter = pythonPath
				}
				settings.ExecStart = settings.Interpreter + " " + targetPath
			}
			if settings.UserName == "" {
				user, err := user.Current()
				handleErr(err)
				settings.UserName = user.Username
			}
			if settings.WorkingDir == "" {
				settings.WorkingDir = workingDir
			}

			// add any extra args to ExecStart
			for i := 1; i < len(args); i++ {
				settings.ExecStart = settings.ExecStart + " " + args[i]
			}

			// Constructing the service file
			tmpl, err := template.ParseFS(TemplateFs, "templates/default-service.tmpl")
			handleErr(err)
			// create a new file
			serviceFilePath := getServicePath(serviceName)
			file, err := os.Create(serviceFilePath)
			handleErr(err)
			// apply the template to the vars map and write the result to file.
			err = tmpl.Execute(file, settings)
			handleErr(err)

			fmt.Printf("Created service file %s\n", serviceFilePath)
			runShell("systemctl", "daemon-reload")
			if !settings.CreateServiceOnly {
				runShell("systemctl", "enable", serviceName)
				runShell("systemctl", "start", serviceName)
			}
		} else if !fileExists(targetPath) {
			// We want to start an existing service
			fmt.Println("Trying to start existing service " + serviceName)
			runShell("systemctl", "start", serviceName)
		} else {
			// We don't know what we want
			fatal(fmt.Sprintf(`
The arguments you provided for 'mp3 start' are ambivalent,
they could either refer to an app in your current working directory %s, or an existing service %s, so mp3 does not know what
it should start.
`, targetPath, serviceName))
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().String("", "", "Pass extra arguments to the script")
	startCmd.Flags().StringVarP(&settings.AppName, "name", "n", "", "Specify an app name")
	startCmd.Flags().StringVarP(&settings.UserName, "user", "u", "", "Specify user who should start the app")
	startCmd.Flags().IntVar(&settings.RestartDelay, "restart-delay", 500, "Delay between automatic restarts")
	startCmd.Flags().BoolVar(&settings.NoAutorestart, "no-autorestart", false, "Do not auto restart app")
	startCmd.Flags().BoolVar(&settings.CreateServiceOnly, "create-only", false, "Create a service file only, do not start or enable")
	startCmd.Flags().StringVar(&settings.Interpreter, "interpreter", "", "Set a specific interpreter to use for executing app, default for .js is /usr/bin/node")
	startCmd.Flags().StringVar(&settings.WorkingDir, "cwd", "", "run target script from path <cwd>")
}
