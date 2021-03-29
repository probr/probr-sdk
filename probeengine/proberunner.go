package probeengine

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/citihub/probr-sdk/audit"
	"github.com/citihub/probr-sdk/config"
	"github.com/cucumber/godog"
)

// ProbeRunner describes the interface that should be implemented to support the execution of tests.
type ProbeRunner interface {
	RunProbe(t *GodogProbe) error
}

// ProbeHandlerFunc describes a callback that should be implemented by test cases in order for ProbeRunner
// to be able to execute the test case.
type ProbeHandlerFunc func(t *GodogProbe) (int, *bytes.Buffer, error)

// GodogProbe encapsulates the specific data that GoDog feature based tests require in order to run.   This
// structure will be passed to the test handler callback.
type GodogProbe struct {
	ProbeDescriptor     *ProbeDescriptor
	ProbeInitializer    func(*godog.TestSuiteContext)
	ScenarioInitializer func(*godog.ScenarioContext)
	FeaturePath         string
	Status              *ProbeStatus `json:"status,omitempty"`
	Results             *bytes.Buffer
}

// RunProbe runs the test case described by the supplied Probe.  It looks in it's test register (the handlers global
// variable) for an entry with the same ProbeDescriptor as the supplied test.  If found, it uses the provided GodogProbe
func (ps *ProbeStore) RunProbe(probe *GodogProbe) (int, error) {

	if probe == nil {
		audit.State.GetProbeLog(probe.ProbeDescriptor.Name).Result = "Internal Error - Probe not found"
		return 2, fmt.Errorf("probe is nil - cannot run test")
	}

	if probe.ProbeDescriptor == nil {
		//update status
		*probe.Status = Error
		audit.State.GetProbeLog(probe.ProbeDescriptor.Name).Result = "Internal Error - Probe descriptor not found"
		return 3, fmt.Errorf("probe descriptor is nil - cannot run test")
	}

	s, o, err := GodogProbeHandler(probe)

	if s == 0 {
		// success
		*probe.Status = CompleteSuccess
	} else {
		// fail
		*probe.Status = CompleteFail
	}

	probe.Results = o // If in-mem output provided, store as Results
	return s, err
}

// RunAllProbes retrieves and executes all probes that have been included
func RunAllProbes(packName string, probes []Probe) (int, *ProbeStore, error) {
	ts := NewProbeStore(packName)

	for _, probe := range probes {
		ts.AddProbe(probe)
	}

	s, err := ts.ExecAllProbes() // Executes all added (queued) tests
	return s, ts, err
}

//GetAllProbeResults maps ProbeStore results to strings
func GetAllProbeResults(ps *ProbeStore) (allResults map[string]string, success bool) {
	allResults = make(map[string]string)
	success = true
	for name := range ps.Probes {
		probeResults, name, err := readProbeResults(ps, name)
		if err != nil {
			allResults[name] = err.Error()
			success = false
		} else {
			allResults[name] = probeResults
		}
	}
	return
}

func readProbeResults(ps *ProbeStore, name string) (probeResults, probeName string, err error) {
	p, err := ps.GetProbe(name)
	if err != nil {
		return
	}
	probeResults = p.Results.String()
	probeName = p.ProbeDescriptor.Name
	return
}

// CleanupTmp is used to dispose of any temp resources used during execution
func CleanupTmp() {
	err := os.RemoveAll(config.Vars.TmpDir())
	if err != nil {
		log.Printf("[ERROR] Error removing tmp folder %v", err)
	}
}
