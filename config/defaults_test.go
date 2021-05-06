package config

import (
	"os"
	"testing"
)

const (
	defaultValuePROBRWRITEDIRECTORY = "probr_output"
	envVarValuePROBRWRITEDIRECTORY  = "ValueFromEnvVar_WriteDirectory"
)

func Test_setFromEnvOrDefaults(t *testing.T) {

	envVarCurrentValuePROBRWRITEDIRECTORY := os.Getenv("PROBR_WRITE_DIRECTORY") // Used to restore env to original state after test
	defer func() {
		os.Setenv("PROBR_WRITE_DIRECTORY", envVarCurrentValuePROBRWRITEDIRECTORY)
	}()

	tests := []struct {
		testName                     string
		envVarValue                  string
		expectedResultWriteDirectory string
	}{
		{
			testName:                     "setFromEnvOrDefaults_GivenEnvVar_ShouldSetConfigVarToEnvVarValue",
			envVarValue:                  envVarValuePROBRWRITEDIRECTORY,
			expectedResultWriteDirectory: envVarValuePROBRWRITEDIRECTORY,
		},
		{
			testName:                     "setFromEnvOrDefaults_WithoutEnvVar_ShouldSetConfigVarToDefaultValue",
			envVarValue:                  "",
			expectedResultWriteDirectory: defaultValuePROBRWRITEDIRECTORY,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {

			os.Setenv("PROBR_WRITE_DIRECTORY", tt.envVarValue)

			vars := &VarOptions{}
			setFromEnvOrDefaults(vars) //This function will modify config object

			//Check WriteDirectory
			if vars.WriteDirectory != tt.expectedResultWriteDirectory {
				t.Errorf("setFromEnvOrDefaults(); PROBR_WRITE_DIRECTORY = %v, Expected: %v", vars.WriteDirectory, tt.expectedResultWriteDirectory)
				return
			}

		})
	}
}
