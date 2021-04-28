package audit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	sdk "github.com/citihub/probr-sdk"
	"github.com/citihub/probr-sdk/config"
	"github.com/citihub/probr-sdk/utils"
)

// SummaryState is a stateful object intended to hold all the high-level info about a probe execution
type SummaryState struct {
	Meta           map[string]interface{}
	Status         string
	ProbesPassed   int
	ProbesFailed   int
	ProbesSkipped  int
	Probes         map[string]*Probe
	WriteDirectory string
}

// SummaryState is a stateful object intended to hold all the high-level info about a probe execution
type limitedSummaryState struct {
	Meta           map[string]interface{}
	Status         string
	ProbesPassed   int
	ProbesFailed   int
	ProbesSkipped  int
	Probes         map[string]*limitedProbe
	WriteDirectory string
}

// NewSummaryState creates a new SummaryState with default values
func NewSummaryState(packName string) (state SummaryState) {
	writeDirectory := filepath.Join(sdk.GlobalConfig.OutputDir(), packName)
	state = SummaryState{
		Probes:         make(map[string]*Probe),
		Meta:           make(map[string]interface{}),
		WriteDirectory: writeDirectory,
	}
	return
}

// TODO: Marshal json, unmarshal into limited obj, then marshal & write/print

// PrintSummary will print the current object state, formatted to JSON
func (s *SummaryState) PrintSummary() {
	log.Printf("Summary: %s", s.summary()) // Summary output should not be handled by log levels
}

// WriteSummary will write the summary to the audit directory
func (s *SummaryState) WriteSummary() {
	path := filepath.Join(config.Vars.GetWriteDirectory(), "summary.json")
	if utils.WriteAllowed(path) {
		ioutil.WriteFile(path, s.summary(), 0755)
	}
}

func (s *SummaryState) summary() []byte {
	var limitedObj limitedSummaryState
	fullJSON := utils.JSON(s)
	err := json.Unmarshal(fullJSON, &limitedObj)
	if err != nil {
		log.Fatalf("[ERROR] Failed to parse summary into JSON: %s", err)
	}
	return utils.JSON(limitedObj)
}

// SetProbrStatus evaluates the current SummaryState state to set the Status
func (s *SummaryState) SetProbrStatus() {
	if s.ProbesPassed > 0 && s.ProbesFailed == 0 {
		s.Status = "Complete - All Probes Completed Successfully"
	} else {
		s.Status = fmt.Sprintf("Complete - %v of %v Probes Failed", s.ProbesFailed, (len(s.Probes) - s.ProbesSkipped))
	}
}

// LogProbeMeta accepts a test name with a key and value to insert to the meta logs for that test. Overwrites key if already present.
func (s *SummaryState) LogProbeMeta(name string, key string, value interface{}) {
	probe := s.GetProbeLog(name)
	probe.Meta[key] = value
	s.Probes[name] = probe
	s.Probes[name].name = name // probe must be able to access its own name, but it is not publicly printed
}

// ProbeComplete takes an probe name and status then updates the summary & probe meta information
func (s *SummaryState) ProbeComplete(name string) {
	p := s.GetProbeLog(name)
	s.completeProbe(p)
	p.Write()
}

// GetProbeLog initializes or returns existing log probe for the provided test name
func (s *SummaryState) GetProbeLog(name string) *Probe {
	// If SummaryState is improperly initialized, a dereference error will occur below.
	// log.Printf("[DEBUG] GetProbeLog(%s) called by: %s->%s->%s", name, utils.CallerName(1), utils.CallerName(2), utils.CallerName(3))
	if s.Probes[name] == nil {
		s.initProbe(name)
	}
	return s.Probes[name]
}

// LogPodName adds pod names to a list for user's debugging purposes
func (s *SummaryState) LogPodName(n string) {
	podNames := s.Meta["names of pods created"].([]string)
	podNames = append(podNames, n)

	s.Meta["names of pods created"] = podNames
}

func (s *SummaryState) initProbe(n string) {
	s.Probes[n] = &Probe{
		name: n,
		Meta: make(map[string]interface{}),
		Path: filepath.Join(config.Vars.AuditDir(), (n + ".json")),
	}
}

func (s *SummaryState) completeProbe(e *Probe) {
	e.countResults()
	if e.Result == "Excluded" {
		e.Meta["audit_path"] = ""
		s.ProbesSkipped = s.ProbesSkipped + 1
	} else if len(e.Scenarios) < 1 {
		e.Result = "No Scenarios Executed"
		e.Meta["audit_path"] = ""
		s.ProbesSkipped = s.ProbesSkipped + 1
	} else if e.ScenariosFailed < 1 {
		e.Result = "Success"
		s.ProbesPassed = s.ProbesPassed + 1
	} else {
		e.Result = "Failed"
		s.ProbesFailed = s.ProbesFailed + 1
	}
}
