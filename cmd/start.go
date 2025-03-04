package cmd

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ShinyTrinkets/overseer"

	quote "github.com/kballard/go-shellquote"
)

// TODO refactor
type Process struct {
	Name  string   `yaml:"name"`
	Cmd   string   `yaml:"cmd"`
	Cwd   string   `yaml:"cwd"`
	Env   []string `yaml:"env"`
	Delay uint	   `yaml:"delay"`
	Retry uint	   `yaml:"retry"`
}


// TODO refactor
type Logger struct {
	Name string
}

func (l *Logger) Debug(msg string, v ...interface{}) {
	msg = "[" + l.Name + "] " + msg
	log.Debugf(msg, v...)
}

func (l *Logger) Info(msg string, v ...interface{}) {
	msg = "[" + l.Name + "] " + msg
	log.Infof(msg, v...)
}

func (l *Logger) Warn(msg string, v ...interface{}) {
	msg = "[" + l.Name + "] " + msg
	log.Warnf(msg, v...)
}

func (l *Logger) Error(msg string, v ...interface{}) {
	msg = "[" + l.Name + "] " + msg
	log.Errorf(msg, v...)
}

func parseProcesses(apps []interface{}) []Process {
	configApps := viper.Get("apps")
	processes := make([]Process, 0)

	for _, app := range configApps.([]interface{}) {
		app := app.(map[string]interface{})

		process := Process{}

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
			log.Warn("Invalid app", "app", app)
		}
	}

	return processes
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the configured commands",
	Long:  "Read the list of commands from the configuration file and start them",
	Run: func(cmd *cobra.Command, args []string) {
		appFile := cmd.Flag("apps").Value.String()

		if appFile != "" {
			log.Info("Reading app file", "file", appFile)
			v := viper.New()

			v.SetConfigFile(appFile)
			v.ReadInConfig()

			apps := v.AllSettings()["apps"]

			if apps != nil {
				viper.MergeConfigMap(map[string]interface{} {"apps": apps})
				log.Info("Apps updated from", "appFile", appFile)
				log.Debug("Config", "config", viper.AllSettings())
			} else {
				log.Warn("No apps found in app file", "appFile", appFile)
			}
		}

		log.Info("Starting commands")

		// TODO refator
		// TODO logger setup
		processes := parseProcesses(viper.Get("apps").([]interface{}))

		log.Debug("Processes", "processes", processes)

		overseer.SetupLogBuilder(func(name string) overseer.Logger {
			return &Logger{Name: name}
		})

		ovr := overseer.NewOverseer()

		for _, process := range processes {
			log.Info("Starting process", "name", process.Name, "cmd", process.Cmd)

			args, err := quote.Split(process.Cmd)
			if err != nil {
				log.Warn("Error parsing command", "cmd", process.Cmd, "error", err)
				continue
			}

			opts := overseer.Options{
				Buffered: false,
				Streaming: true,
				Env: os.Environ(),
			}
			if process.Cwd != "" {
				opts.Dir = process.Cwd
			}
			if len(process.Env) > 0 {
				// TODO check
				opts.Env = append(opts.Env, process.Env...)
			}
			if process.Delay > 0 {
				opts.DelayStart = process.Delay
			}
			if process.Retry > 0 {
				opts.RetryTimes = process.Retry
			}

			p := ovr.Add(process.Name, args[0], args[1:], opts)
			if p == nil {
				log.Warn("Error starting process", "name", process.Name)
				continue
			}
		}

		ovr.SuperviseAll()
		log.Info("All processes stopped")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.PersistentFlags().StringP("apps", "a", "", "App file to read the list of commands from")
}
