# Config Engineering Notes

Internal code to manage config, including Cloud Driver parameters and Test Packs

## Config

Configuration docs are located in the README at the top level of the probr repository.

When creating new config vars, remember to do the following:

1. Add an entry to the struct `ConfigVars` in `internal/config/config.go`
1. Add an entry (matching the config vars struct) to `setEnvOrDefaults` in `internal/config/defaults.go`
1. If appropriate, add logic to `cmd/probr-cli/flags.go`

By following the above steps, you will have accomplished the following:
1. A new variable will be available across the entire probr codebase
1. That variable will have a default value
1. An environment variable can be set to override the default value
1. The env var can be overridden by a provided yaml config file
1. If set, a flag can be used to override the all other values
