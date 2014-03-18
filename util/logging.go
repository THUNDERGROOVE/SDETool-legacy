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
)

func LogInit() {
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic("error opening log file")
	}
	Log = log.New(f, "", log.LstdFlags)
	Info("Log started!")
	if DebugLog {
		//log.SetOutput(os.Stderr)
	}
}

func LErr(s string) {
	if DebugLog {
		term.Red(os.Stderr, "Error: "+s)
		term.White(os.Stderr, "")
		fmt.Println()
	}
	Log.SetPrefix("WARN ")
	Log.Println("Error: ", s)
	Log.SetPrefix("")
}
func Info(s string) {
	if DebugLog {
		term.Cyan(os.Stderr, "Info: "+s)
		term.White(os.Stderr, "")
		fmt.Println()
	}
	Log.SetPrefix("INFO ")
	log.Println(s)
	log.SetPrefix("")
}
