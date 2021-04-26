package opa

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/open-policy-agent/opa/rego"
	//"io/ioutil"
	"log"
)

// Eval evaluates the 'jsonInput' against the specified rego function in the .rego file.
// Returns "true" if the function passes (resource is compliant); false if the function does not pass (resource is non-compliant)
// - regoFilePath = fully qualified filepath to the .rego file
// - regoPackageName = the name of the rego package at the top of the .rego file
// - regoFuncName = the rego function to evaluate against
// - jsonInput = byte representation of the json representation of the resource to eval
func Eval(regoFilePath string, regoPackageName string, regoFuncName string, jsonInput *[]byte) (bool, error) {
	load := make([]string, 1)
	load[0] = regoFilePath

	r := rego.New(
		rego.Query(fmt.Sprintf("x = data.%s.%s", regoPackageName, regoFuncName)),
		rego.Load(load, nil))

	ctx := context.Background()
	query, err := r.PrepareForEval(ctx)
	if err != nil {
		log.Printf("error line 16")
		return false, err
	}
	var input interface{}

	if err := json.Unmarshal(*jsonInput, &input); err != nil {
		log.Printf("error line 32")
		return false, err
	}

	log.Printf("input: %v", input)

	rs, err := query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		log.Printf("error line 38")
		return false, err
	}

	v, ok := rs[0].Bindings["x"].(bool)
	if !ok {
		log.Printf("Did not get bool")
		log.Printf(fmt.Sprintf("rs[0].Bindings[\"x\"] = %v", rs[0].Bindings["x"]))
		return false, fmt.Errorf("Did not get bool")
	}

	return v, nil
}
