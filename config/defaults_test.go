package config

import (
	"os"
	"reflect"
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

func Test_setStringVar(t *testing.T) {

	tests := []struct {
		testName            string
		varName             string
		envVarValue         string
		defaultValue        string
		expectedReturnValue string
	}{
		{
			testName:            "Test that set value is returned when provided",
			varName:             "ENV_VAR_1",
			envVarValue:         envVarValuePROBRWRITEDIRECTORY,
			defaultValue:        defaultValuePROBRWRITEDIRECTORY,
			expectedReturnValue: envVarValuePROBRWRITEDIRECTORY,
		},
		{
			testName:            "Test that default value is returned when no value is provided",
			varName:             "ENV_VAR_2",
			envVarValue:         "",
			defaultValue:        defaultValuePROBRWRITEDIRECTORY,
			expectedReturnValue: defaultValuePROBRWRITEDIRECTORY,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			originalValue := os.Getenv(tt.varName) // Used to restore env to original state after test
			defer func() {
				os.Setenv(tt.varName, originalValue)
			}()

			os.Setenv(tt.varName, tt.envVarValue)

			vars := &VarOptions{}
			value := setStringVar(vars.WriteDirectory, tt.varName, tt.defaultValue)

			if value != tt.expectedReturnValue {
				t.Errorf("setFromEnvOrDefaults(); Return Value = %v, Expected: %v", value, tt.expectedReturnValue)
				return
			}

		})
	}
}

func Test_setStringSliceVar(t *testing.T) {
	defaultValue := []string{"one", "one", "two", "three", "five", "eight"}

	tests := []struct {
		testName            string
		varName             string
		envVarValue         string
		expectedReturnValue []string
	}{
		{
			testName:            "Test that set value is returned when provided",
			varName:             "ENV_VAR_1",
			envVarValue:         "one,two,three,four",
			expectedReturnValue: []string{"one", "two", "three", "four"},
		},
		{
			testName:            "Test that default value is returned when no value is provided",
			varName:             "ENV_VAR_2",
			envVarValue:         "",
			expectedReturnValue: defaultValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			originalValue := os.Getenv(tt.varName) // Used to restore env to original state after test
			defer func() {
				os.Setenv(tt.varName, originalValue)
			}()

			os.Setenv(tt.varName, tt.envVarValue)

			value := setStringSliceVar([]string{}, tt.varName, defaultValue)

			if !reflect.DeepEqual(value, tt.expectedReturnValue) {
				t.Errorf("setFromEnvOrDefaults(); Return Value = %v, Expected: %v", value, tt.expectedReturnValue)
				return
			}

		})
	}
}
