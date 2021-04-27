package probeengine

import (
	"testing"

	"github.com/cucumber/godog"
)

const (
	probeStoreName = "TestProbeStore"
	probeName      = "good_probe"
)

type TestProbe struct {
	name string
}

// Name presents the name of this probe for external reference
func (probe TestProbe) Name() string {
	return probe.name
}

// Path presents the path of these feature files for external reference
func (probe TestProbe) Path() string {
	return ""
}

// ProbeInitialize handles any overall Test Suite initialisation steps.  This is registered with the
// test handler as part of the init() function.
func (probe TestProbe) ProbeInitialize(ctx *godog.TestSuiteContext) {
}

// ScenarioInitialize provides initialization logic before each scenario is executed
func (probe TestProbe) ScenarioInitialize(ctx *godog.ScenarioContext) {
}

func createProbeObj(name string) Probe {
	return &TestProbe{
		name: name,
	}
}

func TestNewProbeStore(t *testing.T) {
	ts := NewProbeStore(probeStoreName)
	if ts == nil {
		t.Logf("Probe store was not initialized")
		t.Fail()
	} else if ts.Probes == nil {
		t.Logf("Probe store was not ready to add probes")
		t.Fail()
	}
}

func TestAddProbe(t *testing.T) {
	ps := NewProbeStore(probeStoreName)
	ps.AddProbe(createProbeObj(probeName))

	// Verify correct conditions succeed
	if ps.Probes[probeName] == nil {
		t.Logf("Probe not added to probe store")
		t.Fail()
	} else if ps.Probes[probeName].Name != probeName {
		t.Logf("Probe name not set properly in test store")
		t.Fail()
	}
}

func TestGetProbe(t *testing.T) {
	ps := NewProbeStore(probeStoreName)
	probe := createProbeObj(probeName)
	ps.AddProbe(probe)

	retrievedProbe, err := ps.GetProbe(probeName)
	if err != nil {
		t.Logf(err.Error())
		t.Fail()
	}
	if retrievedProbe.Name != probe.Name() {
		t.Logf("Retrieved probe does not match added probe")
		t.Fail()
	}
}

// Integration methods:
// TestExecProbe
// TestExecAllProbes
