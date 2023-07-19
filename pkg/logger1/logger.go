package logger1

import (
	"fmt"
	"io"
	"log"
	"strings"
)

type Logger interface {
	Info(msg interface{}, keyvals ...interface{})
	Fatal(msg interface{}, keyvals ...interface{})
	Infof(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	SetOutput(io.Writer)
	SetPrefix(string)
}


var (
	// ErrMissingValue is returned when a key is missing a value.
	ErrMissingValue = fmt.Errorf("missing value")
)

type DefaultLogger struct {
	logger *log.Logger
}

func (l DefaultLogger) SetOutput(w io.Writer) {
	l.logger.SetOutput(w)
}

func (l DefaultLogger) SetPrefix(prefix string) {
	l.logger.SetPrefix(prefix)
}


func (l DefaultLogger) keyvals(keyvals ...interface{}) string {
	if len(keyvals) == 0 {
		return ""
	}

	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, ErrMissingValue)
	}

	var kvs []string
	kvs = append(kvs, " |")
	for i := 0; i <= len(keyvals)/2; i += 2 {
		kvs = append(kvs, fmt.Sprintf("%s=%v", keyvals[i], keyvals[i+1]))
	}

	return strings.Join(kvs, " ")
}

// Fatal implements Logger.
func (l DefaultLogger) Fatal(msg interface{}, keyvals ...interface{}) {
	l.logger.Fatal(fmt.Sprintf("FATAL: %v", msg), l.keyvals(keyvals...))
}

func (l DefaultLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

// Info implements Logger.
func (l DefaultLogger) Info(msg interface{}, keyvals ...interface{}) {
	l.logger.Printf(fmt.Sprintf("INFO: %v", msg), l.keyvals(keyvals...))
}

func (l DefaultLogger) Infof(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

func NewDefaultLogger() Logger {
	return DefaultLogger{
		logger: log.Default(),
	}
}
