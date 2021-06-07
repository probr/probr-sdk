package logging

import (
	"io"
	"log"
	"os"

	hclog "github.com/hashicorp/go-hclog"
)

var (
	activeLogger hclog.Logger
	logger       map[string]hclog.Logger
)

func init() {
	// Initialize default logger
	name := ""
	logger = make(map[string]hclog.Logger)
	logger[name] = newHCLogger(io.Writer(os.Stderr), false)
	writer := logger[name].StandardWriter(&hclog.StandardLoggerOptions{InferLevels: true})
	SetLogWriter(name, writer)
}

// Logger returns the active logger for use in
// statements such as Logger().Info("")
func Logger() hclog.Logger {
	return activeLogger
}

// newHCLogger creates a new hclog.Logger instance with the given name and log level
func newHCLogger(writer io.Writer, jsonFormat bool) hclog.Logger {
	// For level options, reference:
	// https://github.com/hashicorp/go-hclog/blob/master/logger.go#L19
	return hclog.New(&hclog.LoggerOptions{
		Name:       "",
		Level:      hclog.Error,
		Output:     writer,
		JSONFormat: jsonFormat, // TODO: Check env var to determine json format
	})
}

// UpdateLogger ensures that the log package uses
// an hc logger with the provided name and loglevel
// param[0] = Level
// param[1] = Output
// param[2] = JSONFormat
func UpdateLogger(name string, params ...interface{}) {
	if len(params) > 1 {
		// If a writer is specified a new logger will be required
		if len(params) > 2 {
			logger[name] = newHCLogger(params[1].(io.Writer), params[2].(bool))
		} else {
			// If json format is not specified, default to false
			logger[name] = newHCLogger(params[1].(io.Writer), false)
		}
	}
	if logger[name] == nil {
		// Extend base logger for new entry
		logger[name] = logger["default"].Named(name)
	}

	if len(params) > 0 {
		// Set loglevel; Defaults to value in config vars
		logger[name].SetLevel(hclog.LevelFromString(params[0].(string)))
	}
	writer := logger[name].StandardWriter(&hclog.StandardLoggerOptions{InferLevels: true})
	SetLogWriter(name, writer)
}

// SetLogWriter sets the log package to use the provided writer
// If the logWriter was not created using this package,
// Logger() will return nil until UpdateLogger is run
func SetLogWriter(name string, logWriter io.Writer) {
	activeLogger = logger[name]
	log.SetFlags(0)
	log.SetPrefix(name)
	log.SetOutput(logWriter)
}

// GetLogger returns the logger requested by name
func GetLogger(name string) hclog.Logger {
	if logger[name] == nil {
		return nil
	}
	return logger[name]
}
