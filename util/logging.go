package util

import (
	"fmt"
	"github.com/joshlf13/term"
	"log"
	"os"
	"reflect"
)

var (
	Log      *log.Logger
	DebugLog bool
	Color    bool
)

// TypeString was initially made for our logging functions however it's can be
// used all over the codebase
func TypeString(i []interface{}) string {
	s := ""
	for _, v := range i {
		switch k := v.(type) {
		case int:
			s += fmt.Sprintf("%v ", k)
		case string:
			s += fmt.Sprintf("%v ", k)
		case float64:
			s += fmt.Sprintf("%v ", k)
		case SDEType:
			s += fmt.Sprintf("%v ", k.GetName())
		default:
			s += fmt.Sprint("util.TypeString() Does not support type " + reflect.TypeOf(k).String())
		}
	}
	return s
}

// LogInit is called to init the logging portion of the util package.  If you
// try using any of the logging functions before calling this you will get a
// nil pointer exception.
func LogInit() {
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic("error opening log file")
	}
	Color = true
	Log = log.New(f, "", log.Ltime)
	Info("Log started!")
}

// LErr is a function for non-fatal errors.  It will always log to file and
// optionally with -debug print to stdout and also will have nice colored
// output as long as -nocolor is not supplied.
func LErr(i ...interface{}) {
	s := TypeString(i)
	if DebugLog {
		if Color {
			term.Red(os.Stdout, "Error: "+s+"\n")
		} else {
			fmt.Print("Error: " + s)
		}
	}
	Log.SetPrefix("WARN ")
	Log.Println("Error: ", s)
	Log.SetPrefix("")
}

// Info is a function for informing the user what is going on.  It will
// always log to file and optionally with -debug print to stdout and
// also will have nice coloredoutput as long as -nocolor is not supplied.
func Info(i ...interface{}) {
	s := TypeString(i)
	if DebugLog {
		if Color {
			term.Cyan(os.Stdout, "Info: "+s+"\n")

		} else {
			fmt.Print("Info: " + s)
		}
	}
	Log.SetPrefix("INFO ")
	Log.Println(s)
	log.SetPrefix("")
}

// Trace is a function for non-helper functions to call on call.  It will
// always log to file and optionally with -debug print to stdout and
// also will have nice coloredoutput as long as -nocolor is not supplied.
// Don't:
//   Use on primitive, short or otherwise uneeded functions.  An example
//   of ones would be the logging functions
// Do:
//   Use on complicated functions, an example would be most of the sde.go
//   file and any method that uses util.TimeFunction on a defer.
func Trace(i ...interface{}) {
	s := TypeString(i)
	if DebugLog {
		if Color {
			term.Green(os.Stdout, "Trace: "+s+"\n")
		} else {
			fmt.Print("Trace: " + s)
		}
	}
	Log.SetPrefix("TRACE ")
	Log.Println(s)
	log.SetPrefix("")
}
