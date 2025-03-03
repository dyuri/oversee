package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var debug bool

var rootCmd = &cobra.Command{
	Use:   "oversee [command]",
	Short: "Oversee is a tool to execute and monitor commands",
	Long:  "Oversee is a tool to execute and monitor commands",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		err := initViperConfig(cmd)

		return err
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable debug mode")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
}

func initViperConfig(cmd *cobra.Command) error {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")

	v.AddConfigPath("$XDG_CONFIG_HOME/oversee")
	v.AddConfigPath("$HOME/.config/oversee")
	v.AddConfigPath("/etc/oversee")

	// setup log level
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
		log.Debug("No config file found")
		// TODO save the default config file?
	} else {
		// logging
		if v.GetBool("debug") {
			log.SetLevel(log.DebugLevel)
		}

		log.Debug("Config file loaded", "configFile", v.ConfigFileUsed())
	}
	viper.MergeConfigMap(v.AllSettings())
	log.Debug("Config", "config", viper.AllSettings())

	return nil
}

func Execute() {
	log.SetLevel(log.InfoLevel)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
