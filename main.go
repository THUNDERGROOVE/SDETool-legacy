package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/THUNDERGROOVE/SDETool/args"
	"github.com/THUNDERGROOVE/SDETool/config"
	"github.com/THUNDERGROOVE/SDETool/util"
)

var (
	BUILD_DATE string
	VERSION    string
)

func main() {
	defer util.TimeFunction(time.Now(), "main()")
	util.LogInit()
	// If we panic we would like to log where the actual panic came
	// from.  Panics written to console from the default runtime aren't
	// useful to users to report bugs
	defer func() {
		if r := recover(); r != nil {
			util.LErr("Oh noes; We paniced")
			_, file, line, ok := runtime.Caller(0)
			f, _ := filepath.Split(file)
			util.LErr(fmt.Sprintf("%v | %v:%v %v", r, f, line, ok))
		}
	}()

	util.SetCWD()
	args.Init()
	config.LoadConfig()

	// This set is to prevent some nasty nil pointer things if other programs
	// import our packages.  Must be called as soon as possible to prevent errors
	util.VerboseInfo = *args.VerboseInfo
	util.TimeFunc = *args.TimeExecution
	util.DebugLog = *args.Debug
	util.SDEVersion = *args.SDEVersion
	util.Color = util.Inverse(*args.NoColor)

	if util.DebugLog {
		util.Info("Debug logging on")
	}
	// Setup database stuff
	util.CheckFile()
	util.DBInitialize()

	// Change to select switch?
	switch {
	default:
		if *args.SearchFlag == "" && *args.InfoFlag == "" && *args.Damage == "" {
			util.PrintHeader()
			args.PrintHelp()
		}
	case *args.LicenseFlag:
		util.PrintLicense()
	case *args.VersionFlag:
		if VERSION != "" {
			fmt.Println("SDETool version: ", VERSION)
		} else { // Assume was built with `make` or `go build`
			fmt.Println("SDETool version: Dev build")
		}
		if BUILD_DATE != "" {
			fmt.Println("Built on: ", BUILD_DATE)
		} else {
			fmt.Println("No build date specified in binary")
		}
	case *args.DumpTypes:
		fmt.Println("Dumping types to text file :D")
		util.DumpTypes()
	case *args.GetMarketData && *args.InfoFlag == "":
		fmt.Println("The -m(arket) flag requires that you specifiy a type with -i")
	case *args.ForcePanic:
		k := make([]byte, 0)
		print(string(k[99]))
	}
	if *args.SearchFlag != "" {
		fmt.Println("Searching value: '" + *args.SearchFlag + "'")
		s := util.SearchSDEFlag(*args.SearchFlag)
		fmt.Println(s)
	} else if *args.InfoFlag != "" {
		i := util.ResolveInput(*args.InfoFlag)
		t := util.GetSDETypeID(i)
		if *args.ApplyModule != "" { // Apply the module to a suit if we can
		}
		t.PrintInfo()
		if *args.GetMarketData {
			fmt.Println("===== Market Report =====")
			a := t.GetTotalISKSpent()
			t.GetDescription()
			fmt.Println("There has been", util.CommaSep(int64(a)), "ISK spent on", t.GetName())
		}
	} else if *args.Damage != "" {
		t := util.GetSDETypeID(util.ResolveInput(*args.Damage))
		fmt.Println("Getting damage on: " + t.GetName())
		if *args.ComplexModCount == 0 && *args.EnhancedModCount == 0 && *args.BasicModCount == 0 && *args.Prof == 0 {
			t.PrintDamageChart()
			return
		}
		fmt.Println("->", t.GetName(), "would do ", -1, "damage")
		if *args.Prof != 0 {
			fmt.Println("->", "Proficiency level", *args.Prof)
		}
		if *args.ComplexModCount != 0 {
			fmt.Println("->", *args.ComplexModCount, "Complex damage modifiers")
		}
		if *args.EnhancedModCount != 0 {
			fmt.Println("->", *args.EnhancedModCount, "Complex damage modifiers")
		}
		if *args.BasicModCount != 0 {
			fmt.Println("->", *args.BasicModCount, "Complex damage modifiers")
		}
	}
	os.Exit(0)
}
