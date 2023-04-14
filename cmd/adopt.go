/*
Copyright © 2023 Jonas Grünwald
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// adoptCmd represents the adopt command
var adoptCmd = &cobra.Command{
	Use:   "adopt",
	Short: "Adopt a systemd service",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		conn, ctx := connectToSystemd()
		defer conn.Close()
		var targetService = args[0] + ".service"
		var servicePath = findServicePath(targetService, conn, ctx)
		var services = viper.GetStringSlice("AdoptedServices")
		services = append(services, targetService)
		viper.Set("AdoptedServices", services)
		saveConfig()
		fmt.Println(fmt.Sprintf("Adopted service: %s", servicePath))
		fmt.Println(fmt.Sprintf("To unadopt, run mp3 delete %s", args[0]))
	},
}

func init() {
	rootCmd.AddCommand(adoptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// adoptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// adoptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
