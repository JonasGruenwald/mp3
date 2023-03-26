/*
Copyright © 2023 Jonas Grünwald
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a service definition that was created with mp3",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// if adopted service, unadopt the service
		if isAdoptedService(args[0]) {
			fmt.Println(fmt.Sprintf("%s is an adopted service, it will now be unadopted from MP3.", args[0]))
			fmt.Println("If you want to actually delete the service, you will have to do it manually.")
			var services = viper.GetStringSlice("AdoptedServices")
			viper.Set("AdoptedServices", remove(services, getServiceName(args[0])))
			saveConfig()
			return
		}
		var serviceName = getServiceName(args[0])
		// we only want to delete our own services, so we use the path constructor instead of finder
		var targetServicePath = getServicePath(serviceName)

		runSilentWithErr("systemctl", "reset-failed", serviceName)
		runShell("systemctl", "stop", serviceName)
		runShell("systemctl", "disable", serviceName)

		e := os.Remove(targetServicePath)
		if e != nil {
			fatal(e.Error())
		}
		fmt.Println(fmt.Sprintf("Deleted service file %s ", targetServicePath))
		runShell("systemctl", "daemon-reload")
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
