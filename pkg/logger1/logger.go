package logger1

import (
	"log"
)

type Logger interface {
	Info(msg interface{}, keyvals ...interface{})
	Fatal(msg interface{}, keyvals ...interface{})
}

type DefaultLogger struct {
	logger *log.Logger
}

// Fatal implements Logger.
func (l DefaultLogger) Fatal(msg interface{}, keyvals ...interface{}) {
	l.logger.Fatalf(msg.(string), keyvals...)
}

// Info implements Logger.
func (l DefaultLogger) Info(msg interface{}, keyvals ...interface{}) {
	l.logger.Printf(msg.(string), keyvals...)
}

func NewDefaultLogger() Logger {
	return DefaultLogger{
		logger: log.Default(),
	}
}
