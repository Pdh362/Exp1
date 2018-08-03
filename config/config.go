package config

import (
	"encoding/json"
	"github.com/pkg/errors"
	"os"
)

// ------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------
//
// ReadConfig:
//
// Reads configuration from a json file
//
func Read(path string, res interface{}) error {

	// Open config file
	file, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "ReadConfig - Opening configuration file failed:")
	}
	defer file.Close()

	// Read in json settings
	decoder := json.NewDecoder(file)
	err = decoder.Decode(res)
	if err != nil {
		return errors.Wrap(err, "ReadConfig - Config file json decode error:")
	}

	return nil
}
