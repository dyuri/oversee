package config

import (
	"github.com/dyuri/oversee/log"
	"github.com/dyuri/oversee/proc"

	"github.com/spf13/viper"
)

func InitViperConfig(debug bool) error {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")

	v.AddConfigPath("$XDG_CONFIG_HOME/oversee")
	v.AddConfigPath("$HOME/.config/oversee")
	v.AddConfigPath("/etc/oversee")

	// setup log level
	if debug {
		log.SetDebug(true)
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
			log.SetDebug(true)
		}

		log.Debug("Config file loaded", "configFile", v.ConfigFileUsed())
	}
	viper.MergeConfigMap(v.AllSettings())
	log.Debug("Config", "config", viper.AllSettings())

	return nil
}

func ParseProcesses() []proc.Process {
	configApps := viper.Get("apps")
	processes := make([]proc.Process, 0)

	for _, app := range configApps.([]interface{}) {
		app := app.(map[string]interface{})

		process := proc.Process{}

		if name, ok := app["name"].(string); ok {
			process.Name = name
		}
		if cmd, ok := app["cmd"].(string); ok {
			process.Cmd = cmd
		}
		if cwd, ok := app["cwd"].(string); ok {
			process.Cwd = cwd
		}
		if env, ok := app["env"].([]interface{}); ok {
			for _, e := range env {
				process.Env = append(process.Env, e.(string))
			}
		}
		if delay, ok := app["delay"].(int); ok {
			process.Delay = uint(delay)
		}
		if retry, ok := app["retry"].(int); ok {
			process.Retry = uint(retry)
		}

		if process.Name != "" && process.Cmd != "" {
			processes = append(processes, process)
		} else {
			log.Warn("Invalid app: %v", app)
		}
	}

	log.Debug("Processes: %v", processes)
	return processes
}

func UpdateApps(appFile string) {
	if appFile != "" {
		log.Debug("Reading app file: %s", appFile)
		v := viper.New()

		v.SetConfigFile(appFile)
		v.ReadInConfig()

		apps := v.AllSettings()["apps"]

		if apps != nil {
			viper.MergeConfigMap(map[string]interface{}{"apps": apps})
			log.Info("Apps updated from: %s", appFile)
			log.Debug("Config: ", viper.AllSettings())
		} else {
			log.Warn("No apps found in app file: %s", appFile)
		}
	}
}
