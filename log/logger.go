package log

import (
	"github.com/charmbracelet/log"
)

// TODO file based
// TODO log rotation
// TODO Attrs support https://github.com/ShinyTrinkets/meta-logger/blob/master/default.go
// TODO ui support

var DefaultLogger = &Logger{Name: "default"}

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

func (l *Logger) Fatal(msg string, v ...interface{}) {
	msg = "[" + l.Name + "] " + msg
	log.Fatalf(msg, v...)
}

// functions
func SetDebug(debug bool) {
	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

func Debug(msg string, v ...interface{}) {
	DefaultLogger.Debug(msg, v...)
}

func Info(msg string, v ...interface{}) {
	DefaultLogger.Info(msg, v...)
}

func Warn(msg string, v ...interface{}) {
	DefaultLogger.Warn(msg, v...)
}

func Error(msg string, v ...interface{}) {
	DefaultLogger.Error(msg, v...)
}

func Fatal(msg string, v ...interface{}) {
	DefaultLogger.Fatal(msg, v...)
}

