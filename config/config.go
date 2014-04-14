package config

import (
	"encoding/json"
	"github.com/THUNDERGROOVE/SDETool/args"
	"io/ioutil"
	"os"
)

var Conf Config

// Used to set default boolean flags
type Config struct {
	VerboseInfo   bool // If our info should print as much data about a type that we can
	LicenseFlag   bool // Print Licensing information
	VersionFlag   bool // Print current version
	SlowFlag      bool // Don't use optimizations
	TimeExecution bool // Should we time our functions?
	Debug         bool
}

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
func SaveConfig() {
}
