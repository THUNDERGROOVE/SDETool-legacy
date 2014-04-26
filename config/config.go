/*
	The config package is a simple config loaded for SDETool.  It loads a file
	SDETool.config from either %HOME% or $HOME depending on the platform.
*/
package config

import (
	"encoding/json"
	"github.com/THUNDERGROOVE/SDETool/args"
	"io/ioutil"
	"os"
)

var Conf Config

// Config is the data set used the Marshal/Unmarshal our config
type Config struct {
	VerboseInfo   bool // If our info should print as much data about a type that we can
	LicenseFlag   bool // Print Licensing information
	VersionFlag   bool // Print current version
	SlowFlag      bool // Don't use optimizations
	TimeExecution bool // Should we time our functions?
	Debug         bool // Debug man
}

// fcheck is just a small function to check if our config exists and if not it
// make a clean config file
func fcheck() {
	if _, err := os.Stat("SDETool.config"); os.IsNotExist(err) {
		c := Config{false, false, false, false, false, false}
		j, err2 := json.Marshal(c)
		if err2 != nil {
			panic(err2)
		}
		ioutil.WriteFile("SDETool.config", j, 0777)
	}
}

// LoadConfig loads our config from file
func LoadConfig() {
	fcheck()
	f, err1 := ioutil.ReadFile("SDETool.config")
	if err1 != nil {
		panic(err1)
	}
	err2 := json.Unmarshal(f, &Conf)
	if err2 != nil {
		panic(err2)
	}
	if *args.VerboseInfo == false {
		*args.VerboseInfo = Conf.VerboseInfo
	}
	if *args.LicenseFlag == false {
		*args.LicenseFlag = Conf.LicenseFlag
	}
	if *args.VersionFlag == false {
		*args.VersionFlag = Conf.VersionFlag
	}
	if *args.SlowFlag == false {
		*args.SlowFlag = Conf.SlowFlag
	}
	if *args.TimeExecution == false {
		*args.TimeExecution = Conf.TimeExecution
	}
	if *args.Debug == false {
		*args.Debug = Conf.Debug
	}
}

// SaveConfig saves config to disk.  Currently it does nothign as SDETool
// doesn't have in program config editing.
func SaveConfig() {
}
