package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/citihub/probr/utils"
	"gopkg.in/yaml.v2"
)

// GetTags returns Tags, prioritising command line parameter over vars file
func (ctx *Vars) GetTags() string {
	if ctx.Tags == "" {
		ctx.handleTagExclusions() // only process tag exclusions from vars file if not supplied via the command line
	}
	return ctx.Tags
}

// SetTags will parse the tags specified in Vars.Tags
func (ctx *Vars) SetTags(tags map[string][]string) {
	configTags := strings.Split(ctx.GetTags(), ",")
	for _, configTag := range configTags {
		for _, tag := range tags[configTag] {
			configTags = append(configTags, "@"+tag)
		}
	}
	ctx.Tags = strings.Join(configTags, ",")
}

// Handle tag exclusions provided via the config vars file
func (ctx *Vars) handleTagExclusions() {
	for _, tag := range ctx.TagExclusions {
		if ctx.Tags == "" {
			ctx.Tags = "~@" + tag
		} else {
			ctx.Tags = fmt.Sprintf("%s && ~@%s", ctx.Tags, tag)
		}
	}
}

// Init will override config.Vars with the content retrieved from a filepath
func Init(configPath string) (Vars, error) {
	config, err := NewConfig(configPath)

	if err != nil {
		//log.Printf("[ERROR] %v", err)
		return config, err
	}
	setFromEnvOrDefaults(&config)            // Set any values not retrieved from file
	SetLogFilter(config.LogLevel, os.Stderr) // Set the minimum log level obtained from Vars
	//log.Printf("[DEBUG] Config initialized by %s", utils.CallerName(1))

	return config, nil
}

// NewConfig overrides the current config.Vars values
func NewConfig(c string) (Vars, error) {
	// Create config structure
	config := Vars{}
	if c == "" {
		return config, nil // No file path provided, return empty config
	}
	err := ValidateConfigPath(c)
	if err != nil {
		return config, err
	}
	// Open config file
	file, err := os.Open(c)
	if err != nil {
		return config, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return config, err
	}

	return config, nil
}

// ValidateConfigPath simply ensures the file exists
func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

// LogConfigState will write the config file to the write directory
func (ctx *Vars) LogConfigState() {
	json, _ := json.MarshalIndent(ctx, "", "  ")
	//log.Printf("[INFO] Config State: %s", json)
	path := filepath.Join(ctx.GetWriteDirectory(), "config.json")
	if ctx.WriteConfig == "true" && utils.WriteAllowed(path, ctx.Overwrite()) {
		data := []byte(json)
		ioutil.WriteFile(path, data, 0644)
		//log.Printf("[NOTICE] Config State written to file %s", path)
	}
}

// TmpDir creates and returns -tmp- directory within WriteDirectory
func (ctx *Vars) TmpDir() string {
	tmpDir := filepath.Join(ctx.GetWriteDirectory(), "tmp")
	_ = os.Mkdir(tmpDir, 0755) // Creates if not already existing
	return tmpDir
}

// Overwrite returns the string value of the OverwriteHistoricalAudits in bool format
func (ctx *Vars) Overwrite() bool {
	value, err := strconv.ParseBool(ctx.OverwriteHistoricalAudits)
	if err != nil {
		//log.Printf("[ERROR] Could not parse value '%s' for OverwriteHistoricalAudits %s", ctx.OverwriteHistoricalAudits, err)
		return false
	}
	return value
}

// AuditDir creates and returns -audit- directory within WriteDirectory
func (ctx *Vars) AuditDir() string {
	auditDir := filepath.Join(ctx.GetWriteDirectory(), "audit")
	_ = os.Mkdir(auditDir, 0755) // Creates if not already existing
	return auditDir
}

// CucumberDir creates and returns -cucumber- directory within WriteDirectory
func (ctx *Vars) CucumberDir() string {
	cucumberDir := filepath.Join(ctx.GetWriteDirectory(), "cucumber")
	_ = os.Mkdir(cucumberDir, 0755) // Creates if not already existing
	return cucumberDir
}

// GetWriteDirectory creates and returns the output folder specified in settings
func (ctx *Vars) GetWriteDirectory() string {
	_ = os.Mkdir(ctx.WriteDirectory, 0755) // Creates if not already existing
	return ctx.WriteDirectory
}
