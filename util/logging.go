package util

import (
	//"io/ioutil"
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

func LErr(i ...interface{}) {
	s := ""
	for _, v := range i {
		switch k := v.(type) {
		case int:
			s += fmt.Sprintf("%v ", k)
		case string:
			s += fmt.Sprintf("%v ", k)
		case float64:
			s += fmt.Sprintf("%v ", k)
		default:
			LErr("util.LErr() Does not support type " + reflect.TypeOf(v).String())
		}
	}
	if DebugLog {
		if Color {
			term.Red(os.Stdout, "Error: "+s)
		} else {
			fmt.Print("Error: " + s)
		}
		fmt.Println()
	}
	Log.SetPrefix("WARN ")
	Log.Println("Error: ", s)
	Log.SetPrefix("")
}
func Info(i ...interface{}) {
	s := ""
	for _, v := range i {
		switch k := v.(type) {
		case int:
			s += fmt.Sprintf("%v ", k)
		case string:
			s += fmt.Sprintf("%v ", k)
		case float64:
			s += fmt.Sprintf("%v ", k)
		default:
			LErr("util.Info() Does not support type " + reflect.TypeOf(v).String())
		}
	}
	if DebugLog {
		if Color {
			term.Cyan(os.Stdout, "Info: "+s)

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
			term.Green(os.Stdout, "Trace: "+s)

		} else {
			fmt.Print("Trace: " + s)
		}
		fmt.Println()
	}
	Log.SetPrefix("TRACE ")
	Log.Println(s)
	log.SetPrefix("")
}
