package cliflags

import (
	"fmt"
	"os"
	"testing"
)

func TestFlags_ExecuteHandlers(t *testing.T) {
	testArgs := make(map[string]string)
	testFuncOutput := make(map[string]string)

	testArgs["runTestFunc"] = "set the test func value"
	testArgs["runTestFunc2"] = "set another test func value"
	testArgs["runTestFunc3"] = "be really redundant"

	for key, value := range testArgs {
		os.Args = append(os.Args, fmt.Sprintf("-%s=%s", key, value))
	}

	var flags Flags
	for key := range testArgs {
		flags.NewStringFlag(key, "no description", func(value *string) {
			testFuncOutput[key] = *value
		})
	}

	flags.ExecuteHandlers()

	for _, flag := range flags.PreParsedFlags {
		name := flag.(StringFlag).Name
		value := *flag.(StringFlag).Value
		if testArgs[name] != value {
			t.Errorf("Expected value for '%s' to be '%s', but found '%v'", name, value, testArgs[name])
		}
	}

}

func TestFlags_NewStringFlag(t *testing.T) {
	var testFuncOutput string

	type args struct {
		name    string
		usage   string
		handler stringHandlerFunc
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test that flag is created",
			args: args{
				name:  "testFunc1",
				usage: "testFunc1 usage",
				handler: func(value *string) {
					testFuncOutput = *value
				},
			},
		},
		{
			name: "Test that flag is still created",
			args: args{
				name:  "testFunc2",
				usage: "testFunc2 usage",
				handler: func(value *string) {
					testFuncOutput = *value
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var flags Flags
			expected := "someString/" + tt.args.name
			os.Args = append(os.Args, fmt.Sprintf("-%s=%s", tt.args.name, expected))
			flags.NewStringFlag(tt.args.name, tt.args.usage, tt.args.handler)
			flags.ExecuteHandlers()
			if testFuncOutput == tt.args.name {
				t.Errorf("Expected test args to contain '%s', but found %s", expected, testFuncOutput)
			}
		})
	}
}
