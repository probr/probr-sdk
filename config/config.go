package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/probr/probr-sdk/config/setter"
	"github.com/probr/probr-sdk/utils"
)

// GlobalConfig ...
var GlobalConfig GlobalOpts

// Init ...
func (ctx *GlobalOpts) Init() {
	ctx.StartTime = time.Now()
	if ctx.VarsFile != "" {
		ctx.decode()
	} else {
		log.Printf("[DEBUG] No vars file provided, unexpected behavior may occur")
	}

	ctx.setEnvAndDefaults()

	log.Printf("[DEBUG] Config initialized by %s", utils.CallerName(1))
	ctx.LogConfigState()
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

// LogConfigState ...
func (ctx *GlobalOpts) LogConfigState() {
	json, _ := json.MarshalIndent(ctx, "", "  ")
	log.Printf("[INFO] Config State: %s", json)
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
	setter.SetVar(&ctx.InstallDir, "PROBR_INSTALL_DIR", filepath.Join(home, "probr"))

	setter.SetVar(&ctx.TmpDir, "PROBR_TMP", filepath.Join(ctx.InstallDir, "tmp"))
	setter.SetVar(&ctx.BinariesDir, "PROBR_BIN", filepath.Join(ctx.InstallDir, "bin"))
	setter.SetVar(&ctx.WriteDirectory, "PROBR_WRITE_DIRECTORY", ctx.OutputDir())
	setter.SetVar(&ctx.LogLevel, "PROBR_LOG_LEVEL", "DEBUG")
	setter.SetVar(&ctx.GodogResultsFormat, "PROBR_RESULTS_FORMAT", "cucumber")
}

// ParseTags takes two lists of tags and parses them into a cucumber tag string
// Tags may start with '@' or '~@' respectively, but it is not required
func ParseTags(inclusions, exclusions []string) string {
	var tags []string
	if len(inclusions) > 0 {
		tags = append(tags, parseInclusions(inclusions))
	}
	if len(exclusions) > 0 {
		tags = append(tags, parseExclusions(exclusions))
	}
	// If only one is provided, this joiner won't be used
	return strings.Join(tags, " && ")
}

func parseInclusions(inclusions []string) string {
	inclusions = prependTags(inclusions, "@")
	return strings.Join(inclusions, ",")
}

func parseExclusions(exclusions []string) string {
	exclusions = prependTags(exclusions, "~@")
	return strings.Join(exclusions, " && ")
}

func prependTags(tags []string, prefix string) []string {
	for i, value := range tags {
		// ensure value is not empty, then force it to begin with prefix
		// prefix value is expected to not be found anywhere but as the prefix
		if len(value) > 0 && !strings.Contains(value, prefix) {
			tags[i] = prefix + value
		}
	}
	return tags
}

// SetTmpDir sets the location that temporary files will be written to
func SetTmpDir(path string) {
	GlobalConfig.TmpDir = path
	err := os.MkdirAll(path, 0755)
	if err == nil {
		log.Printf("[DEBUG] Created temporary directory: %v", err)
	} else {
		log.Printf("[ERROR] Failed to create temporary directory: %v", err)
	}
}

// CleanupTmp is used to dispose of any temp resources used during execution
func (ctx *GlobalOpts) CleanupTmp() {
	err := os.RemoveAll(ctx.TmpDir)
	if err != nil {
		log.Printf("[ERROR] Failed to remove temporary directory %v", err)
	}
}

// OutputDir parses a filepath based on GlobalOpts.InstallDir and the datetime this was initialized
func (ctx *GlobalOpts) OutputDir() string {
	year, month, day := ctx.StartTime.Date()
	hour, min, sec := ctx.StartTime.Clock()
	yearMonthDay := fmt.Sprintf("%04d%v%02d", year, month, day)
	timestamp := fmt.Sprintf("%02d%02d%02d", hour, min, sec)
	execName := utils.GetExecutableName()

	base := filepath.Join(ctx.InstallDir, "output", yearMonthDay, timestamp, execName)
	prepareOutputDirectories(base)
	return base
}

func prepareOutputDirectories(base string) {
	dirs := []string{"audit", "cucumber"}
	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(base, dir), 0755)
		if err != nil {
			log.Print(utils.ReformatError(err.Error()))
		} else {
			log.Printf("[DEBUG] Directory is ready for use: %s", dir)
		}
	}
}
