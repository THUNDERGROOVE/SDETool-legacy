/*
	args.go is used so all packages can access our arguments
*/

package args

import (
	"flag"
)

var (
	// Flag variables
	SearchFlag    *string // SearchFlag is used to provide a string to search for types
	InfoFlag      *string // InfoFlag is used to provide an int to display info about a type
	VerboseInfo   *bool   // If our info should print as much data about a type that we can
	LicenseFlag   *bool   // Print Licensing information
	VersionFlag   *bool   // Print current version
	SlowFlag      *bool   // Don't use optimizations
	TimeExecution *bool   // Should we time our functions?
	Clean         *bool   // Cleans cache and database files
	DumpTypes     *bool   // Dumps all types to a file for use with category.go
	GetMarketData *bool   // Flag used if getting market data with -i

	// Damage calculations
	Damage           *string // Damage is used to provide a TypeID to calculate damage of a weapon
	ComplexModCount  *int    // ComplexModCount is used to calculate how many Complex mods to use
	EnhancedModCount *int    // EnhancedModCount is used to calculate how many Enhanced mods to use
	BasicModCount    *int    // BasicModCount is used to calculate how many Basic mods to use
	Prof             *int    // Prof is how many levels of proficiency used when calculating damage

	// Server related
	RunServer *bool // Runs a server for hosting the web version of SDETool
	Port      *int  // What port to listen on?
)

func init() {
	// Flags
	SearchFlag = flag.String("s", "", "Search for TypeIDs")
	InfoFlag = flag.String("i", "", "Get info with a TypeID, typeName or mDisplayName")
	VerboseInfo = flag.Bool("vi", false, "Prints all attributes when used with -i")
	LicenseFlag = flag.Bool("l", false, "Prints license information.")
	VersionFlag = flag.Bool("version", false, "Prints the SDETool version")
	SlowFlag = flag.Bool("slow", false, "Forces the use of unoptimized functions")
	TimeExecution = flag.Bool("time", false, "Times the execution of functions that may take a decent amount of time")
	Clean = flag.Bool("clean", false, "Cleans all database and cache files")
	DumpTypes = flag.Bool("dump", false, "Dumps all types to a file for use with the category package")
	GetMarketData = flag.Bool("m", false, "Gets market data on item, used with -i. Sorry CCP if I'm pounding your APIs ;P")

	// Damage and mod counts
	Damage = flag.String("d", "", "Get damage calculations, takes a TypeID")
	ComplexModCount = flag.Int("c", 0, "Amount of complex damage mods, used with -d")
	EnhancedModCount = flag.Int("e", 0, "Amount of enhanced damage mods, used with -d")
	BasicModCount = flag.Int("b", 0, "Amount of enhanced damage mods, used with -d")
	Prof = flag.Int("p", 0, "Prof level, used with -d")

	// Server related
	RunServer = flag.Bool("server", false, "Runs a server for hosting the web version of SDETool")
	Port = flag.Int("port", 80, "Port used for -server")
	flag.Parse()
}
