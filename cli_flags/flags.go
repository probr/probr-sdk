package cliflags

import (
	"flag"
	"log"
	"os"

	"github.com/citihub/probr-sdk/config"
	"github.com/citihub/probr-sdk/utils"
)

type stringHandlerFunc func(value *string)
type boolHandlerFunc func(value *bool)

// Flag allows for different value types to be handled
type Flag interface {
	executeHandler()
}

// Flags allows for the gathering all flag definitions and executing them together
type Flags struct {
	PreParsedFlags []Flag
}

// StringFlag holds the user-provided value for the flag, and the function to be run within executeHandler
type StringFlag struct {
	Name    string
	Handler stringHandlerFunc
	Value   *string
}

// BoolFlag holds the user-provided value for the flag, and the function to be run within executeHandler
type BoolFlag struct {
	Name    string
	Handler boolHandlerFunc
	Value   *bool
}

func (f StringFlag) executeHandler() {
	f.Handler(f.Value)
}

func (f BoolFlag) executeHandler() {
	f.Handler(f.Value)
}

// ExecuteHandlers executes the logic for any flags that are provided via `./probr (--<FLAG>)`
func (flags *Flags) ExecuteHandlers() {
	flag.Parse()
	for _, f := range flags.PreParsedFlags {
		f.executeHandler()
	}
}

// NewStringFlag creates a new flag that accepts string values
func (flags *Flags) NewStringFlag(name string, usage string, handler stringHandlerFunc) {
	f := StringFlag{
		Name:    name,
		Handler: handler,
		Value:   new(string),
	}
	flag.StringVar(f.Value, name, "", usage)
	flags.PreParsedFlags = append(flags.PreParsedFlags, f)
}

// NewBoolFlag creates a new flag that accepts bool values
func (flags *Flags) NewBoolFlag(name string, usage string, handler boolHandlerFunc) {
	f := BoolFlag{
		Handler: handler,
		Value:   new(bool),
	}
	flag.BoolVar(f.Value, name, false, usage)
	flags.PreParsedFlags = append(flags.PreParsedFlags, f)
}

// VarsFileHandler initializes configuration with VarsFile overriding env vars & defaults
func VarsFileHandler(v *string) {
	value := *v
	err := config.Init(value)
	if err != nil {
		log.Fatalf("[ERROR] error returned from config.Init: %v", err)
	} else if len(value) > 0 {
		config.Vars.VarsFile = value
		log.Printf("[INFO] Config read from file '%v', but may still be overridden by CLI flags.", value)
	} else {
		log.Printf("[NOTICE] No configuration variables file specified. Using environment variabls and defaults only.")
	}
}

// WriteDirHandler changes the root output directory
func WriteDirHandler(v *string) {
	value := *v
	if len(value) > 0 {
		log.Printf("[NOTICE] Output Directory has been overridden via command line")
		config.Vars.WriteDirectory = value
	}
}

// LoglevelHandler validates provided value is a known loglevel and sets loglevel accordingly
// TODO: does this still work with our new logger?
func LoglevelHandler(v *string) {
	value := *v
	if len(value) > 0 {
		levels := []string{"DEBUG", "INFO", "NOTICE", "WARN", "ERROR"}
		_, found := utils.FindString(levels, value)
		if !found {
			log.Fatalf("[ERROR] Unknown loglevel specified: '%s'. Must be one of %v", value, levels)
		} else {
			config.Vars.LogLevel = value
			config.SetLogFilter(config.Vars.LogLevel, os.Stderr)
		}
	}
}

// ResultsformatHandler parses a flag and sets the godog output type
func ResultsformatHandler(v *string) {
	value := *v
	if len(value) > 0 {
		options := []string{"cucumber", "events", "junit", "pretty", "progress"}
		_, found := utils.FindString(options, value)
		if !found {
			log.Fatalf("[ERROR] Unknown resultsformat specified: '%s'. Must be one of %v", value, options)
		} else {
			config.Vars.ResultsFormat = value
			config.SetLogFilter(config.Vars.ResultsFormat, os.Stderr)
		}
	}
}

// TagsHandler parses a flag and sets the godog/cucumber tags
func TagsHandler(v *string) {
	value := *v
	if len(value) > 0 {
		config.Vars.Tags = value
		log.Printf("[INFO] tags have been added via command line.")
	}
}

// TODO: we might not need this anymore
func isFlagPassed(flagName string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == flagName {
			found = true
		}
	})
	return found
}
