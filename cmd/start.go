package cmd

import (
	"github.com/dyuri/oversee/config"
	"github.com/dyuri/oversee/proc"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the configured commands",
	Long:  "Read the list of commands from the configuration file and start them",
	Run: func(cmd *cobra.Command, args []string) {
		processes := config.ParseProcesses()
		proc.InitProcesses(processes)

		proc.SuperviseAll()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
