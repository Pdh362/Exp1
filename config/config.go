package config

import (
	"encoding/json"
	"flag"
	"github.com/pkg/errors"
	"os"
)

var Mode string
var WPort int
var MPort int

// ------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------
//
// ReadConfig:
//
// Reads configuration from a json file
// These hold global config information, used in either mode that rarely are changed.
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

	// Command-line flag parsing
	flag.StringVar(&Mode, "mode", "watcher", "choose from : watcher=monitor folder, master=serve results.")
	flag.IntVar(&WPort, "Watch port", 8080, "Network port for a watcher to talk to the master")
	flag.IntVar(&MPort, "Master port", 80, "Network port for master to serve results to.")

	flag.Parse()

	return nil
}
