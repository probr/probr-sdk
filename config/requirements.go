package config

// Requirements is used to dictate the required config vars for each service pack
var Requirements = map[string][]string{
	"Storage":    {"Provider"},
	"APIM":       {"Provider"},
	"Kubernetes": {"AuthorisedContainerRegistry", "UnauthorisedContainerRegistry"},
}
