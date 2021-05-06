package config

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/citihub/probr-sdk/utils"
)

// GlobalConfig ...
var GlobalConfig GlobalOpts

// CloudProviders config options
type CloudProviders struct {
	Azure Azure `yaml:"Azure"`
}

// Azure config options that may be required by any service pack
type Azure struct {
	Excluded         string `yaml:"Excluded"`
	TenantID         string `yaml:"TenantID"`
	SubscriptionID   string `yaml:"SubscriptionID"`
	ClientID         string `yaml:"ClientID"`
	ClientSecret     string `yaml:"ClientSecret"`
	ResourceGroup    string `yaml:"ResourceGroup"`
	ResourceLocation string `yaml:"ResourceLocation"`
	ManagementGroup  string `yaml:"ManagementGroup"`
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
	SetVar(&ctx.InstallDir, "PROBR_RESULTS_FORMAT", filepath.Join(home, "probr"))

	SetVar(&ctx.TmpDir, "PROBR_RESULTS_FORMAT", "cucumber")
	SetVar(&ctx.WriteDirectory, "PROBR_WRITE_DIRECTORY", "probr_output")
	SetVar(&ctx.LogLevel, "PROBR_LOG_LEVEL", "DEBUG")
	SetVar(&ctx.GodogResultsFormat, "PROBR_RESULTS_FORMAT", "cucumber")

	SetVar(&ctx.CloudProviders.Azure.TenantID, "PROBR_AZURE_TENANT_ID", "")
	SetVar(&ctx.CloudProviders.Azure.SubscriptionID, "PROBR_AZURE_SUBSCRIPTION_ID", "")
	SetVar(&ctx.CloudProviders.Azure.ClientID, "PROBR_AZURE_CLIENT_ID", "")
	SetVar(&ctx.CloudProviders.Azure.ClientSecret, "PROBR_AZURE_CLIENT_SECRET", "")
	SetVar(&ctx.CloudProviders.Azure.ResourceGroup, "PROBR_AZURE_RESOURCE_GROUP", "")
	SetVar(&ctx.CloudProviders.Azure.ResourceLocation, "PROBR_AZURE_RESOURCE_LOCATION", "")
}
