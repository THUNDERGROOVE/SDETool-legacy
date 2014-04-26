/*
	args is a package used to keep all command line parsing out util so our
	code can be reused
*/
package args

import (
	"flag"
)

var (
	// Flag variables
	SearchFlag       *string // SearchFlag is used to provide a string to search for types
	InfoFlag         *string // InfoFlag is used to provide an int to display info about a type
	Damage           *string // Damage is used to provide a TypeID to calculate damage of a weapon
	SDEVersion       *string // Version of the SDE to force
	ApplyModule      *string
	VerboseInfo      *bool // If our info should print as much data about a type that we can
	LicenseFlag      *bool // Print Licensing information
	VersionFlag      *bool // Print current version
	SlowFlag         *bool // Don't use optimizations
	TimeExecution    *bool // Should we time our functions?
	Clean            *bool // Cleans cache and database files
	DumpTypes        *bool // Dumps all types to a file for use with category.go
	GetMarketData    *bool // Flag used if getting market data with -i
	RunServer        *bool // Runs a server for hosting the web version of SDETool
	Debug            *bool
	ForcePanic       *bool
	Quiet            *bool
	Uninstall        *bool
	NoColor          *bool
	NoSkills         *bool // Calculate stats without skill bonuses
	ComplexModCount  *int  // ComplexModCount is used to calculate how many Complex mods to use
	EnhancedModCount *int  // EnhancedModCount is used to calculate how many Enhanced mods to use
	BasicModCount    *int  // BasicModCount is used to calculate how many Basic mods to use
	Prof             *int  // Prof is how many levels of proficiency used when calculating damage
	Port             *int  // What port to listen on?
)

// Init handles parsing command line flags
func Init() {
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
	ApplyModule = flag.String("m", "", "Used with -i to apply a module to a dropsuit")

	GetMarketData = flag.Bool("market", false, "Gets market data on item, used with -i. Sorry CCP if I'm pounding your APIs ;P")
	Debug = flag.Bool("debug", false, "Debug? Debug!")
	ForcePanic = flag.Bool("fp", false, "Forces a panic, debug uses")
	Quiet = flag.Bool("quiet", false, "Used with flags like uninstall where you want it to produce no output, ask for input or block in any sort of way")
	Uninstall = flag.Bool("uninstall", false, "Uninstalls SDETool if install via makefile or manually in your PATH variable")
	NoColor = flag.Bool("nocolor", false, "Used to disable color.  Usefull for >'ing and |'ing")
	NoSkills = flag.Bool("ns", false, "Used to prevent SDETool from applying skill bonuses")
	// Damage and mod counts
	Damage = flag.String("d", "", "Get damage calculations, takes a TypeID")
	SDEVersion = flag.String("sv", "1.8", "Version of the SDE to use, in the form of '1.7' or '1.8'")
	ComplexModCount = flag.Int("c", 0, "Amount of complex damage mods, used with -d")
	EnhancedModCount = flag.Int("e", 0, "Amount of enhanced damage mods, used with -d")
	BasicModCount = flag.Int("b", 0, "Amount of enhanced damage mods, used with -d")
	Prof = flag.Int("p", 0, "Prof level, used with -d")

	// Server related
	RunServer = flag.Bool("server", false, "Runs a server for hosting the web version of SDETool")
	Port = flag.Int("port", 80, "Port used for -server")
	flag.Parse()
}
