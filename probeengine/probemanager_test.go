package probeengine

import (
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

// TODO: Add tests for ProbeStore
