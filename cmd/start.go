package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strings"
	"text/template"
)

type ServiceSettings struct {
	AppName       string
	RestartDelay  int
	NoAutorestart bool
	Interpreter   string
	UserName      string
	WorkingDir    string
	ExecStart     string
}

var settings = ServiceSettings{}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start an application with pmu",
	Long: `Supply the name of an application entry point in your current working directory,
and pmu will create a daemonized service for it for you.
Pass the name of an existing pmu-created service to start it

Examples:

# Spin up new service:
pmu start bashscript.sh
pmu start python-app.py
pmu start ./binary-file -- --port 1520

# Start existing service
pmu start my-app

`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		/**
		We need to figure out if the user wants to spin up a service, or start an existing service.
		We will check which option is possible and decide based on that
		*/
		target := args[0]
		workingDir, err := os.Getwd()
		if err != nil {
			fatal("Could not get working directory")
		}
		// path to potential new app to create service for
		var targetPath = path.Join(workingDir, target)
		// path to potential existing service to start
		var serviceName = fmt.Sprintf("%s%s.service", serviceNamePrefix, target)
		var targetServicePath = path.Join(systemCtlUnitDir, serviceName)

		if settings.AppName != "" || (fileExists(targetPath) && !fileExists(targetServicePath)) {
			// We want to create a new service
			// Supplement info for service file
			if settings.AppName != "" {
				serviceName = fmt.Sprintf("%s%s.service", serviceNamePrefix, settings.AppName)
			}
			settings.ExecStart = targetPath
			// If we are dealing with a script file, we need to add the interpreter
			if strings.HasSuffix(target, ".js") {
				if settings.Interpreter == "" {
					settings.Interpreter = "/usr/bin/node"
				}
				settings.ExecStart = settings.Interpreter + " " + targetPath
			} else if strings.HasSuffix(target, ".py") {
				if settings.Interpreter == "" {
					settings.Interpreter = "/usr/bin/python3"
				}
				settings.ExecStart = settings.Interpreter + " " + targetPath
			}
			if settings.UserName == "" {
				user, err := user.Current()
				if err != nil {
					fmt.Println(err.Error())
					fatal("Can't get user!")
				}
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
			tmpl, err := template.ParseFiles("service-template.tmpl")
			if err != nil {
				fatal("Error parsing template file for " + serviceName)
			}

			// create a new file
			file, err := os.Create(path.Join(systemCtlUnitDir, serviceName))
			if err != nil {
				fatal("Error touching service file for " + serviceName)
			}
			defer func(file *os.File) {
				err := file.Close()
				if err != nil {

				}
			}(file)

			// apply the template to the vars map and write the result to file.
			err = tmpl.Execute(file, settings)
			if err != nil {
				fatal("Error generating service file for " + serviceName)
			}

			fmt.Printf("Created new service %s", serviceName)

			// execute template and write to file
		} else if !fileExists(targetPath) {
			// We want to start an existing service
			cmd := exec.Command("systemctl", "start", serviceName)
			output, err := cmd.CombinedOutput()
			fmt.Printf("%s\n", output)
			if err != nil {
				fatal(fmt.Sprintf("pmu encountered an error trying to start the service %s", serviceName))
			}
		} else {
			// We don't know what we want
			fatal(fmt.Sprintf(`
The arguments you provided for 'pmu start' are ambivalent,
they could either refer to an app in your current working directory %s, or an existing service %s, so pmu does not know what
it should start.
`, targetPath, serviceName))
		}

		fmt.Println("Start called !")
		fmt.Println(fmt.Sprintf("Delay: %dms", settings.RestartDelay))
		fmt.Println(fmt.Sprintf("Name: %s", settings.AppName))
		fmt.Println(fmt.Sprintf("NoAutorestart: %t", settings.NoAutorestart))
		fmt.Println(fmt.Sprintf("Args: %s", args))
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().String("", "", "Pass extra arguments to the script")
	startCmd.Flags().StringVarP(&settings.AppName, "name", "n", "", "Specify an app name")
	startCmd.Flags().StringVarP(&settings.UserName, "user", "u", "", "Specify user who should start the app")
	startCmd.Flags().IntVar(&settings.RestartDelay, "restart-delay", 500, "Delay between automatic restarts")
	startCmd.Flags().BoolVar(&settings.NoAutorestart, "no-autorestart", false, "Do not auto restart app")
	startCmd.Flags().StringVar(&settings.Interpreter, "interpreter", "", "Set a specific interpreter to use for executing app, default for .js is /usr/bin/node")
	startCmd.Flags().StringVar(&settings.WorkingDir, "cwd", "", "run target script from path <cwd>")
}
