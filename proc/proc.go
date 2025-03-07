package proc

import (
	"os"

	"github.com/dyuri/oversee/log"

	"github.com/ShinyTrinkets/overseer"

	quote "github.com/kballard/go-shellquote"
)

type Process struct {
	Name  string   `yaml:"name"`
	Cmd   string   `yaml:"cmd"`
	Cwd   string   `yaml:"cwd"`
	Env   []string `yaml:"env"`
	Delay uint	   `yaml:"delay"`
	Retry uint	   `yaml:"retry"`
}

var ovr *overseer.Overseer

func init() {
	overseer.SetupLogBuilder(func(name string) overseer.Logger {
		return &log.Logger{Name: name}
	})

	ovr = overseer.NewOverseer()
}

func InitProcesses(processes []Process) {
	for _, process := range processes {
		args, err := quote.Split(process.Cmd)
		if err != nil {
			log.Warn("Error parsing command [%s]: %s", process.Cmd, err)
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
			log.Warn("Error starting process: %s", process.Name)
			continue
		}
	}
}

func GetOverseer() *overseer.Overseer {
	return ovr
}

func SuperviseAll() {
	ovr.SuperviseAll()
}
