package config

import (
	"testing"
)

//
// Tests
//

func TestNewConfig(t *testing.T) {
	// Just use a default config, no file-read for now
	config, err := NewConfig("")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	switch interface{}(config).(type) {
	case VarsObject:
	default:
		t.Log("NewConfig did not create a VarsObject object")
		t.Fail()
	}
}

// TestOverwrite ...
func TestOverwrite(t *testing.T) {
	vars, _ := NewConfig("")
	vars.OverwriteHistoricalAudits = "true"
	if vars.Overwrite() != true {
		t.Errorf("Overwrite() should return a bool of 'true'")
	}

	vars.OverwriteHistoricalAudits = "false"
	if vars.Overwrite() != false {
		t.Errorf("Overwrite() should return a bool of 'false'")
	}
}
