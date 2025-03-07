package cmd

import (
	"github.com/dyuri/oversee/log"
	"github.com/dyuri/oversee/config"
	"github.com/dyuri/oversee/proc"
	"github.com/dyuri/oversee/ui"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var debug bool // TODO remove?

var rootCmd = &cobra.Command{
	Use:   "oversee [command]",
	Short: "Oversee is a tool to execute and monitor commands",
	Long:  "Oversee is a tool to execute and monitor commands",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		err := config.InitViperConfig(debug)

		appFile := cmd.Flag("apps").Value.String()
		config.UpdateApps(appFile)

		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		processes := config.ParseProcesses()
		proc.InitProcesses(processes)

		ui.StartUI()
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable debug mode")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	rootCmd.PersistentFlags().StringP("apps", "a", "", "App file to read the list of commands from")
}

func Execute() {
	log.SetDebug(false)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("%v", err)
	}
}
