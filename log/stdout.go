package log

import (
	"log"

	"go.uber.org/zap"
)

// StdoutLogger is logger implementation that logs to stdout
type StdoutLogger struct {
	logger *zap.Logger
}

// NewLogger return new StdoutLogger instance with configuration based on the config
func NewLogger(verbose bool) *StdoutLogger {
	var logger *zap.Logger

	switch verbose {
	case true:
		devLogger, err := zap.NewDevelopment()
		if err != nil {
			log.Panic("can't initialize logger")
		}
		logger = devLogger
	default:
		prodLogger, err := zap.NewProduction()
		if err != nil {
			log.Panic("can't initialize logger")
		}
		logger = prodLogger
	}
	stdoutLogger := &StdoutLogger{
		logger: logger,
	}
	return stdoutLogger
}

// LogInfo creates information level log
func (l *StdoutLogger) LogInfo(i string) {
	defer l.logger.Sync()
	l.logger.Info("INFO", zap.String("msg", i))
}

// LogError creates error level log along with error details
func (l *StdoutLogger) LogError(msg string, err error) {
	defer l.logger.Sync()
	l.logger.Error(msg, zap.Error(err))
}

// LogFatal creates fatal level log which also terminates execution
func (l *StdoutLogger) LogFatal(msg string, err error) {
	defer l.logger.Sync()
	l.logger.Fatal(msg, zap.Error(err))
}
