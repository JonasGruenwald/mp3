package cmd

import (
	"embed"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mp3",
	Short: "A process management utility (not really)",
	Long: `MP3 is a convenient CLI for creating and managing systemd services.
Its syntax is based on the popular node.js process manager 'pm2'.

Start an app with mp3 start app.js, then check its status with mp3 status.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !printStatus() {
			err := cmd.Help()
			handleErr(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(templateFs embed.FS) {
	TemplateFs = templateFs
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mp3.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
