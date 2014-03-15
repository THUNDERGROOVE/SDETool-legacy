package main

import (
	"flag"
	"fmt"
	"github.com/THUNDERGROOVE/SDETool/args"
	"github.com/THUNDERGROOVE/SDETool/server"
	"github.com/THUNDERGROOVE/SDETool/util"
	"os"
)

const (
	// Constants
	SDEFile = "dustSDE.db"                                                           // Name for the SDE database file to be used
	SDEUrl  = "http://cdn1.eveonline.com/community/DUST_SDE/Uprising_1.7_674383.zip" // URL to download the SDE
	Version = 0.2                                                                    // Sounds right :P
)

func main() {
	checkfile()
	args.Init()
	util.DBInitialize()
	// Change to select switch?
	if *args.LicenseFlag {
		util.PrintLicense()
	} else if *args.VersionFlag {
		fmt.Println("SDETool version", Version)
	} else if *args.RunServer {
		server.RunServer()
	} else if *args.SearchFlag != "" {
		fmt.Println("Searching value: '" + *args.SearchFlag + "'")
		s := util.SearchSDEFlag(*args.SearchFlag)
		fmt.Println(s)
	} else if *args.InfoFlag != "" {
		i := util.ResolveInput(*args.InfoFlag)
		t := util.GetSDETypeID(i)
		t.PrintInfo()
		if *args.GetMarketData {
			fmt.Println("===== Market Report =====")
			a := t.GetTotalISKSpent()
			fmt.Println("There has been", a, "ISK spent on", t.GetName())
		}
	} else if *args.Damage != "" {
		t := util.GetSDETypeID(util.ResolveInput(*args.Damage))
		fmt.Println("Getting damage on: " + t.GetName())
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
		os.Remove(SDEFile)
		os.Remove(SDEFile + ".zip")
	} else if *args.DumpTypes {
		fmt.Println("Dumping types to text file :D")
		util.DumpTypes()
	} else if *args.GetMarketData && *args.InfoFlag == "" {
		fmt.Println("The -m(arket) flag requires that you specifiy a type with -i")
	} else {
		flag.PrintDefaults()
	}
}
