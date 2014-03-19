package util

import (
	//"io/ioutil"
	"fmt"
	"github.com/joshlf13/term"
	"log"
	"os"
)

var (
	Log      *log.Logger
	DebugLog bool
	Color    bool
)

func LogInit() {
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic("error opening log file")
	}
	Color = true
	Log = log.New(f, "", log.Ltime)
	Info("Log started!")
	if DebugLog {
		//log.SetOutput(os.Stderr)
	}
}

func LErr(s string) {
	if DebugLog {
		if Color {
			term.Red(os.Stderr, "Error: "+s)
		} else {
			fmt.Print("Error: " + s)
		}
		fmt.Println()
	}
	Log.SetPrefix("WARN ")
	Log.Println("Error: ", s)
	Log.SetPrefix("")
}
func Info(s string) {
	if DebugLog {
		if Color {
			term.Cyan(os.Stderr, "Info: "+s)

		} else {
			fmt.Print("Info: " + s)
		}
		fmt.Println()
	}
	Log.SetPrefix("INFO ")
	Log.Println(s)
	log.SetPrefix("")
}
func Trace(s string) {
	if DebugLog {
		if Color {
			term.Green(os.Stderr, "Trace: "+s)

		} else {
			fmt.Print("Trace: " + s)
		}
		fmt.Println()
	}
	Log.SetPrefix("TRACE ")
	Log.Println(s)
	log.SetPrefix("")
}
