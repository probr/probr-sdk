package cliflags

import (
	"fmt"
	"os"
	"testing"
)

var testFuncOutput []string

func testFunc(value *string) {
	testFuncOutput = append(testFuncOutput, *value)
}

func TestFlags_ExecuteHandlers(t *testing.T) {
	testArgs := make(map[string]string)
	testArgs["runTestFunc"] = "set the test func value"
	testArgs["runTestFunc2"] = "set another test func value"
	testArgs["runTestFunc3"] = "be really redundant in this test"

	for key, value := range testArgs {
		os.Args = append(os.Args, fmt.Sprintf("-%s=%s", key, value))
	}

	var flags Flags
	for key, _ := range testArgs {
		flags.NewStringFlag(key, "no description", testFunc)
	}
	flags.ExecuteHandlers()

	i := 0
	for _, value := range testArgs {
		if testFuncOutput[i] != value {
			t.Errorf("Expected testArgs to contain '%s', but did not find it", value)
			return
		}
		i = i + 1
	}
}
