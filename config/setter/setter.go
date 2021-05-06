package setter

import (
	"log"
	"os"
	"strings"
)

// SetVar fetches the env var or sets the default value as needed for the specified field from VarOptions
func SetVar(field interface{}, varName string, defaultValue interface{}) {
	switch v := field.(type) {
	case *string:
		*field.(*string) = setStringVar(*field.(*string), varName, defaultValue.(string))
	case *[]string:
		*field.(*[]string) = setStringSliceVar(*field.(*[]string), varName, defaultValue.([]string))
	default:
		log.Fatalf("Unexpected value type provided for '%v', should be %T", varName, v)
	}
}

func setStringVar(value string, varName string, defaultValue string) string {
	if value == "" { // if field was empty, get value from env var
		value = os.Getenv(varName)
	}
	if value == "" { // if still empty, use default value provided
		value = defaultValue
	}
	return value
}

func setStringSliceVar(value []string, varName string, defaultValue []string) []string {
	if len(value) == 0 { // if field was empty, get value from env var
		t := os.Getenv(varName) // for []string, env var should be comma separated values
		if len(t) > 0 {
			value = append(value, strings.Split(t, ",")...)
		}
	}
	if len(value) == 0 { // if still empty, use default value provided
		value = defaultValue
	}
	return value
}
