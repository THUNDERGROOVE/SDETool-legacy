package main

import (
	"flag"
	"fmt"
)

const (
	SDEFile = "dustSDE.db"
	SDEUrl  = "http://cdn1.eveonline.com/community/DUST_SDE/Uprising_1.7_674383.zip"
)

var (
	// Single commands
	SearchFlag  *string
	InfoFlag    *int
	VerboseInfo *bool

	// Damage calculations
	Damage           *int
	ComplexModCount  *int
	EnhancedModCount *int
	BasicModCount    *int
	Prof             *int
)

func init() {
	// Flags
	SearchFlag = flag.String("s", "", "Search for TypeIDs")
	InfoFlag = flag.Int("i", 0, "Get info with TypeID")
	VerboseInfo = flag.Bool("vi", false, "Prints all attributes when used with -i")

	Damage = flag.Int("d", 0, "Get damage calculations, takes a TypeID")
	ComplexModCount = flag.Int("c", 0, "Amount of complex damage mods, used with -d")
	EnhancedModCount = flag.Int("e", 0, "Amount of enhanced damage mods, used with -d")
	BasicModCount = flag.Int("b", 0, "Amount of enhanced damage mods, used with -d")
	Prof = flag.Int("p", 0, "Prof level, used with -d")

	flag.Parse()
}
func main() {
	checkfile()
	DBInitialize()
	if *SearchFlag != "" {
		fmt.Println("Searching value: '" + *SearchFlag + "'")
		k := GetSDEWhereNameContains(*SearchFlag)
		for _, c := range k {
			fmt.Println(c.TypeID, "| "+c.GetName())
		}
	} else if *InfoFlag != 0 {
		t := GetSDETypeID(*InfoFlag)
		t.PrintInfo()
	} else if *Damage != 0 {
		t := GetSDETypeID(*Damage)
		fmt.Println("Getting damage on: " + t.GetName())
		d := t.GetRawDamage(*Prof, *ComplexModCount, *EnhancedModCount, *BasicModCount)
		fmt.Println("->", t.GetName(), "would do ", d, "damage")
		if *Prof != 0 {
			fmt.Println("->", "Proficiency level", *Prof)
		}
		if *ComplexModCount != 0 {
			fmt.Println("->", *ComplexModCount, "Complex damage modifiers")
		}
		if *EnhancedModCount != 0 {
			fmt.Println("->", *EnhancedModCount, "Complex damage modifiers")
		}
		if *BasicModCount != 0 {
			fmt.Println("->", *BasicModCount, "Complex damage modifiers")
		}

	} else {
		flag.PrintDefaults()
	}
}
