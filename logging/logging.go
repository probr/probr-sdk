package logging

import (
	"io"
	"log"
	"os"

	hclog "github.com/hashicorp/go-hclog"
)

var (
	// logger is the global hclog logger
	logger hclog.Logger

	// logWriter is a global writer for logs, to be used with the std log package
	logWriter io.Writer
)

func init() {
	logger = newHCLogger("")
	logWriter = logger.StandardWriter(&hclog.StandardLoggerOptions{InferLevels: true})

	// set up the default std library logger to use our output
	log.SetFlags(0)
	log.SetPrefix("")
	log.SetOutput(logWriter)
}

// newHCLogger returns a new hclog.Logger instance with the given name
func newHCLogger(name string) hclog.Logger {
	logOutput := io.Writer(os.Stderr)

	logLevel := hclog.Trace // TODO: Add logic to check env var and determine log level

	jsonFormat := true // TODO: Add logic to check env var and determine json format

	// TODO: Add logic to check env var and determine if should log to file, then prepare file "refistersink"

	return hclog.New(&hclog.LoggerOptions{
		Name:       name,
		Level:      logLevel,
		Output:     logOutput,
		JSONFormat: jsonFormat,
	})
}

// ProbrLogger returns the default global hclog logger
func ProbrLogger() hclog.Logger {
	return logger
}

// ProbrLoggerOutput returns the log writer for default global hclog logger
func ProbrLoggerOutput() io.Writer {
	return logWriter
}
