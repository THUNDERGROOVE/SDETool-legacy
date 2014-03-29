package main

import (
	"flag"
	"fmt"
	"github.com/THUNDERGROOVE/SDETool/args"
	"github.com/THUNDERGROOVE/SDETool/server"
	"github.com/THUNDERGROOVE/SDETool/util"
	"io/ioutil"
	"os"
	"time"
)

const (
	Version = 0.3
)

var (
	BuildDate string
)

func main() {
	defer util.TimeFunction(time.Now(), "main()")
	defer func() {
		if r := recover(); r != nil {
			util.DebugLog = true // Make it print and log always
			fname := fmt.Sprintf("%s%s", time.Now().String(), "SDETool.panic")
			util.LErr("Panic detected writing to file " + fname)
			err := ioutil.WriteFile(fname, []byte(fmt.Sprintf("%v", r)), 0777)
			if err != nil {
				util.LErr("SUPER ERROR!?!?!?! Couldn't write panic log :/")
			}
			fmt.Println("Report this pls")
			os.Exit(1)
		}
	}()
	util.SetCWD()
	util.LogInit()
	args.Init()

	// This set is to prevent some nasty nil pointer things if other programs import our packages.  Must be called as soon as possible to prevent errors
	util.VerboseInfo = *args.VerboseInfo
	util.TimeFunc = *args.TimeExecution
	util.DebugLog = *args.Debug
	util.SDEVersion = *args.SDEVersion

	util.Info("Debug logging on")
	util.CheckFile()
	util.DBInitialize()
	// Change to select switch?
	if *args.LicenseFlag {
		util.PrintLicense()
	} else if *args.VersionFlag {
		fmt.Println("SDETool version", Version)
		if BuildDate != "" {
			fmt.Println("Built on: ", BuildDate)
		} else {
			fmt.Println("No build date specified in binary")
		}
	} else if *args.RunServer {
		*args.Debug = true
		server.RunServer()
	} else if *args.SearchFlag != "" {
		fmt.Println("Searching value: '" + *args.SearchFlag + "'")
		s := util.SearchSDEFlag(*args.SearchFlag)
		fmt.Println(s)
	} else if *args.InfoFlag != "" {
		i := util.ResolveInput(*args.InfoFlag)
		t := util.GetSDETypeID(i)
		t.ApplySkillsToType()
		if *args.ApplyModule != "" { // Apply the module to a suit if we can
			g := util.ResolveInput(*args.ApplyModule)
			m := util.GetSDETypeID(g)
			t.ApplyModuleToSuit(m)
		}
		t.PrintInfo()
		if *args.GetMarketData {
			fmt.Println("===== Market Report =====")
			a := t.GetTotalISKSpent()
			fmt.Println("There has been", a, "ISK spent on", t.GetName())
		}
	} else if *args.Damage != "" {
		t := util.GetSDETypeID(util.ResolveInput(*args.Damage))
		fmt.Println("Getting damage on: " + t.GetName())
		if *args.ComplexModCount == 0 && *args.EnhancedModCount == 0 && *args.BasicModCount == 0 && *args.Prof == 0 {
			t.PrintDamageChart()
			return
		}
		d := t.GetRawDamage(*args.Prof, *args.ComplexModCount, *args.EnhancedModCount, *args.BasicModCount)
		fmt.Println("->", t.GetName(), "would do ", d, "damage")
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
	} else if *args.Clean {
		fmt.Println("Cleaning SDETool directory")
		os.Remove(util.SDEFile)
		os.Remove(util.SDEFile + ".zip")
	} else if *args.DumpTypes {
		fmt.Println("Dumping types to text file :D")
		util.DumpTypes()
	} else if *args.GetMarketData && *args.InfoFlag == "" {
		fmt.Println("The -m(arket) flag requires that you specifiy a type with -i")
	} else if *args.ForcePanic {
		util.ForcePanic()
	} else {
		util.PrintHeader()
		flag.PrintDefaults()
	}
}
