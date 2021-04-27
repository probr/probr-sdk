package sdk

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// GlobalOpts provides configurable options that will be used throughout the SDK
type GlobalOpts struct {
	InstallDir         string
	TmpDir             string
	GodogResultsFormat string
	StartTime          time.Time
}

// GlobalConfig allows certain values to be configured at runtime for the entire SDK
var GlobalConfig GlobalOpts

// Sets default values; assumes Probr is installed to user's home directory
// Override by changing values within sdk.GlobalConfig
func init() {
	home, _ := os.UserHomeDir()

	GlobalConfig = GlobalOpts{
		InstallDir:         filepath.Join(home, "probr"),
		GodogResultsFormat: "cucumber",
		StartTime:          time.Now(),
	}
	SetTmpDir(filepath.Join(home, "probr", "tmp")) // TODO: this needs error handling
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
func (gc *GlobalOpts) CleanupTmp() {
	err := os.RemoveAll(gc.TmpDir)
	if err != nil {
		log.Printf("[ERROR] Failed to remove temporary directory %v", err)
	}
}

// OutputDir parses a filepath based on GlobalOpts.InstallDir and the datetime this was initialized
func (gc *GlobalOpts) OutputDir() string {
	year, month, day := gc.StartTime.Date()
	hour, min, sec := gc.StartTime.Clock()
	yearMonthDay := fmt.Sprintf("%d%d%d", year, month, day)
	timestamp := fmt.Sprintf("%d%d%d", hour, min, sec)

	return filepath.Join(gc.InstallDir, "output", yearMonthDay, timestamp)
}
