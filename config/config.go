package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/briandowns/spinner"
	"github.com/citihub/probr-sdk/utils"
	"gopkg.in/yaml.v2"
)

// Vars is a singleton instance of VarOptions
var Vars VarOptions

// Spinner holds the current state of the CLI spinner
var Spinner *spinner.Spinner

// GetTags returns Tags, prioritising command line parameter over vars file
func (ctx *VarOptions) GetTags() string {
	if ctx.Tags == "" {
		ctx.handleTagExclusions() // only process tag exclusions from vars file if not supplied via the command line
	}
	return ctx.Tags
}

// SetTags will parse the tags specified in Vars.Tags
func (ctx *VarOptions) SetTags(tags map[string][]string) {
	configTags := strings.Split(ctx.GetTags(), ",")
	for _, configTag := range configTags {
		for _, tag := range tags[configTag] {
			configTags = append(configTags, "@"+tag)
		}
	}
	ctx.Tags = strings.Join(configTags, ",")
}

// Handle tag exclusions provided via the config vars file
func (ctx *VarOptions) handleTagExclusions() {
	for _, tag := range ctx.TagExclusions {
		if ctx.Tags == "" {
			ctx.Tags = "~@" + tag
		} else {
			ctx.Tags = fmt.Sprintf("%s && ~@%s", ctx.Tags, tag)
		}
	}
}

// Init will override config.GlobalConfig with the content retrieved from a filepath
func Init(configPath string) error {
	config, err := NewConfig(configPath)

	if err != nil {
		log.Printf("[ERROR] %v", err)
		return err
	}
	config.Meta = Vars.Meta // Persist any existing Meta data
	Vars = config
	setFromEnvOrDefaults(&Vars) // Set any values not retrieved from file

	SetLogFilter(Vars.LogLevel, os.Stderr) // Set the minimum log level obtained from Vars
	log.Printf("[DEBUG] Config initialized by %s", utils.CallerName(1))

	Vars.handleConfigFileExclusions()

	return nil
}

// NewConfig overrides the current config.GlobalConfig values
func NewConfig(c string) (VarOptions, error) {
	// Create config structure
	config := VarOptions{}
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
func (ctx *VarOptions) LogConfigState() {
	json, _ := json.MarshalIndent(Vars, "", "  ")
	log.Printf("[INFO] Config State: %s", json)
	path := filepath.Join(ctx.GetWriteDirectory(), "config.json")
	if ctx.WriteConfig == "true" && utils.WriteAllowed(path) {
		data := []byte(json)
		ioutil.WriteFile(path, data, 0644)
		log.Printf("[NOTICE] Config State written to file %s", path)
	}
}

// Overwrite returns the string value of the OverwriteHistoricalAudits in bool format
func (ctx *VarOptions) Overwrite() bool {
	value, err := strconv.ParseBool(ctx.OverwriteHistoricalAudits)
	if err != nil {
		log.Printf("[ERROR] Could not parse value '%s' for OverwriteHistoricalAudits %s", ctx.OverwriteHistoricalAudits, err)
		return false
	}
	return value
}

// AuditDir creates and returns -audit- directory within WriteDirectory
func (ctx *VarOptions) AuditDir() string {
	auditDir := filepath.Join(ctx.GetWriteDirectory(), "audit")
	_ = os.MkdirAll(auditDir, 0755) // Creates if not already existing
	return auditDir
}

// CucumberDir creates and returns -cucumber- directory within WriteDirectory
func (ctx *VarOptions) CucumberDir() string {
	cucumberDir := filepath.Join(ctx.GetWriteDirectory(), "cucumber")
	_ = os.MkdirAll(cucumberDir, 0755) // Creates if not already existing
	return cucumberDir
}

// GetWriteDirectory creates and returns the output folder specified in settings + executable name
func (ctx *VarOptions) GetWriteDirectory() string {

	execName := utils.GetExecutableName()

	// Add executable name to WriteDirectory
	outputDir := filepath.Join(ctx.WriteDirectory, execName)
	_ = os.Mkdir(outputDir, 0755) // Creates if not already existing

	return outputDir
}

func (ctx *VarOptions) handleConfigFileExclusions() {
	ctx.handleProbeExclusions("kubernetes", ctx.ServicePacks.Kubernetes.Probes) //TODO: logic specific to service packs should be handled outside of SDK
	ctx.handleProbeExclusions("storage", ctx.ServicePacks.Storage.Probes)       //TODO: logic specific to service packs should be handled outside of SDK
}

func (ctx *VarOptions) handleProbeExclusions(packName string, probes []Probe) {
	for _, probe := range probes {
		if probe.IsExcluded() {
			ctx.addExclusion(fmt.Sprintf("probes/%s/%s", packName, probe.Name))
		} else {
			for _, scenario := range probe.Scenarios {
				if scenario.IsExcluded() {
					ctx.addExclusion(fmt.Sprintf("probes/%s/%s/%s", packName, probe.Name, scenario.Name))
				}
			}
		}
	}
}

func (ctx *VarOptions) addExclusion(tag string) {
	if len(ctx.Tags) > 0 {
		ctx.Tags = ctx.Tags + " && "
	}
	ctx.Tags = fmt.Sprintf("%s~@%s", ctx.Tags, tag)
}

// IsExcluded will log and return exclusion configuration
func (k Kubernetes) IsExcluded() bool {
	return validatePackRequirements("Kubernetes", k)
}

// IsExcluded will log and return exclusion configuration
func (s Storage) IsExcluded() bool {
	return validatePackRequirements("Storage", s)
}

// IsExcluded will log and return exclusion configuration
func (a APIM) IsExcluded() bool {
	return validatePackRequirements("APIM", a)
}

// IsExcluded will log and return exclusion configuration
func (p Probe) IsExcluded() bool {
	return p.Excluded != ""
}

// IsExcluded will log and return exclusion configuration
func (s Scenario) IsExcluded() bool {
	return s.Excluded != ""
}

func validatePackRequirements(name string, object interface{}) bool {
	// reflect for dynamic type querying
	storage := reflect.Indirect(reflect.ValueOf(object))

	for i, requirement := range Requirements[name] {
		if storage.FieldByName(requirement).String() == "" {
			if Vars.Meta.RunOnly == "" || strings.EqualFold(Vars.Meta.RunOnly, name) {
				// Warn if the pack may have been expected to run
				log.Printf("[WARN] Ignoring %s service pack due to required var '%s' not being present.", name, Requirements[name][i])
			}
			return true
		}
	}
	if Vars.Meta.RunOnly != "" && strings.EqualFold(Vars.Meta.RunOnly, name) {
		// If another pack is specified as RunOnly, this should be excluded
		log.Printf("[NOTICE] Ignoring %s service pack due to %s being specified by 'probr run <SERVICE-PACK-NAME>'", name, Vars.Meta.RunOnly)
		return true
	}
	log.Printf("[NOTICE] %s service pack included.", name)
	return false
}

// GetPacks returns a list of pack names (as specified by internal/config/requirements.go)
func GetPacks() (keys []string) {
	for value := range Requirements {
		keys = append(keys, value)
	}
	return keys
}
