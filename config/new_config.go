package config

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/citihub/probr-sdk/config/setter"
	azureconfig "github.com/citihub/probr-sdk/providers/azure/config"
	"github.com/citihub/probr-sdk/utils"
)

// GlobalConfig ...
var GlobalConfig GlobalOpts

// CloudProviders config options
type CloudProviders struct {
	Azure azureconfig.Azure `yaml:"Azure"`
}

// GlobalOpts provides configurable options that will be used throughout the SDK
type GlobalOpts struct {
	StartTime          time.Time
	VarsFile           string
	InstallDir         string         `yaml:"InstallDir"`
	TmpDir             string         `yaml:"TmpDir"`
	GodogResultsFormat string         `yaml:"GodogResultsFormat"`
	CloudProviders     CloudProviders `yaml:"CloudProviders"`
	WriteDirectory     string         `yaml:"WriteDirectory"`
	LogLevel           string         `yaml:"LogLevel"`
	TagExclusions      []string       `yaml:"TagExclusions"`
	TagInclusions      []string       `yaml:"TagInclusions"`
	WriteConfig        string         `yaml:"WriteConfig"`
}

// Init ...
func (ctx *GlobalOpts) Init() {
	if ctx.VarsFile != "" {
		ctx.decode()
	} else {
		log.Printf("[DEBUG] No vars file provided, unexpected behavior may occur")
	}

	ctx.setEnvAndDefaults()

	log.Printf("[DEBUG] Config initialized by %s", utils.CallerName(1))
}

// decode uses an SDK helper to create a YAML file decoder,
// parse the file to an object, then extracts the values from
// ServicePacks.Kubernetes into this context
func (ctx *GlobalOpts) decode() (err error) {
	configDecoder, file, err := NewConfigDecoder(ctx.VarsFile)
	if err != nil {
		return
	}
	err = configDecoder.Decode(&ctx)
	file.Close()
	return
}

// SetTmpDir sets the location that temporary files will be written to
func (ctx *GlobalOpts) SetTmpDir(path string) {
	ctx.TmpDir = path
	err := os.MkdirAll(path, 0755)
	if err == nil {
		log.Printf("[DEBUG] Created temporary directory: %v", err)
	} else {
		log.Printf("[ERROR] Failed to create temporary directory: %v", err)
	}
}

// setEnvOrDefaults will set value from os.Getenv and default to the specified value
func (ctx *GlobalOpts) setEnvAndDefaults() {
	// Notes on SetVar's values:
	// 1. Pointer to local object; will be overwritten by env or default if empty
	// 2. Name of env var to check
	// 3. Default value to set if flags, vars file, and env have not provided a value

	home, _ := os.UserHomeDir()
	setter.SetVar(&ctx.InstallDir, "PROBR_RESULTS_FORMAT", filepath.Join(home, "probr"))

	setter.SetVar(&ctx.TmpDir, "PROBR_RESULTS_FORMAT", "cucumber")
	setter.SetVar(&ctx.WriteDirectory, "PROBR_WRITE_DIRECTORY", "probr_output")
	setter.SetVar(&ctx.LogLevel, "PROBR_LOG_LEVEL", "DEBUG")
	setter.SetVar(&ctx.GodogResultsFormat, "PROBR_RESULTS_FORMAT", "cucumber")

	ctx.CloudProviders.Azure.SetEnvAndDefaults()
}
