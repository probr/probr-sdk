package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// NewConfigDecoder reads the provided into a new yaml decoder and leaves the file open
func NewConfigDecoder(path string) (decoder *yaml.Decoder, file *os.File, err error) {
	file, err = openConfigPath(path)
	if err != nil {
		return // err if path can not be validated
	}
	decoder = yaml.NewDecoder(file) // Init new YAML decode, leave file open
	return
}

func openConfigPath(path string) (file *os.File, err error) {
	s, err := os.Stat(path)
	if err == nil && s.IsDir() {
		err = fmt.Errorf("'%s' is a directory; expected a normal file", path)
		return
	}
	file, err = os.Open(path)
	if err != nil {
		return
	}
	return
}
