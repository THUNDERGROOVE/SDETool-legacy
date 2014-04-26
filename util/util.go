/*
	The util package is for constants, package variables and functions that belong no where else.
	Used to consolidate most of the codebase out of the main package and allow other devs
	to use my code if they want.
*/
package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	TimeFunc bool
)

// Used to measure execuation time of a function, use with defer
func TimeFunction(start time.Time, name string) {
	elapsed := time.Since(start)
	if TimeFunc {
		Trace(fmt.Sprintf("%s took %s", name, elapsed))
	}
}

const (
	sqlliteDriver = `SDETool uses an SQLite3 driver to deal with the SDE database from https://github.com/mattn/go-sqlite3, see http://mattn.mit-license.org/2012 for licensing`
)

// SDEToolInit should initialize all parts of SDETool, mainly the logging and DB stuffs helpful for making external tools
func SDEToolInit() {
	SetCWD()
	LogInit()
	CheckFile()
	DBInitialize()
}

// SDEToolInitLocal is like SDEToolInit except it doesn't setup an application
// directory and uses the current working directory for configs and database
// files
func SDEToolInitLocal() {
	LogInit()
	CheckFile()
	DBInitialize()
}

// xspaces returns a string of spaces with a length of x
func xspaces(x int) string {
	var t string
	for i := 0; i < x; i++ {
		t = t + " "
	}
	return t
}

// PrintLicense is a simple function to just print our licensing for everything
func PrintLicense() {
	fmt.Println("SDETool is under the MIT license.  See LICENSE for more info")
	fmt.Println(sqlliteDriver)
}

// longestLen returns the length of the longest string in list
func longestLen(list map[string]string) int {
	l := 0
	for v := range list {
		if len(v) > l {
			l = len(v)
		}
	}
	return l
}

// longestLen returns the length of the longest string in list
func longestLenS(list []string) int {
	l := 0
	for _, v := range list {
		if len(v) > l {
			l = len(v)
		}
	}
	return l
}

// ResolveInput takes in a string and returns a TypeID.
// Can take a TypeID, name, display name.  If there are multiple matches
// we will return the closest match
func ResolveInput(s string) int {
	// Check if we have a TypeID
	b, err := regexp.MatchString("^[0-9]{1,6}$", s)
	if err != nil {
		fmt.Println(err.Error())
	}
	if b {
		i, err := strconv.Atoi(s)
		if err != nil {
			fmt.Println(err.Error())
			return 0
		}
		return i
	}
	// Check if we have a name
	if strings.Contains(s, "_") {
		return GetTypeIDByName(s)
	} else {
		return GetTypeIDByDName(s)
	}
	return 0
}

// Helper function designed for PrintInfo() to print a weapon/slottype only
// if not zero.
func printNotZero(name string, value int) {
	if value != 0 {
		fmt.Println(name, value)
	}
}
func printFNotZero(name string, value float64) {
	if value != float64(0) {
		fmt.Println(name, value)
	}
}

func returnNotZero(name string, value int) string {
	if value != 0 {
		return name + " " + strconv.Itoa(value) + "\n"
	} else {
		return ""
	}
}

func cleanTypeName(typeName string) string {
	a := strings.Split(typeName, " ")
	t := strings.Join(a, "") // Remove spaces
	t = strings.Join(strings.Split(t, "-"), "")
	t = strings.Join(strings.Split(t, "'"), "")
	t = strings.Join(strings.Split(t, "/"), "")
	return t
}

// Dumps types into a text file for use in category.go
// DISCLAIMER:  This function is VERY memory intensive.
// Expect it to eat ALL of your RAM and stop to a crawl
// while your OS pages the data :P
func DumpTypes() {
	f := ""
	var err error
	rows, err := db.Query("SELECT * FROM CatmaTypes")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for rows.Next() {
		var typeID int
		var typeName string
		err = rows.Scan(&typeID, &typeName)
		fmt.Println("Proccessing " + typeName)
		s := GetSDETypeIDFast(typeID)
		f += cleanTypeName(s.GetName()) + " = " + strconv.Itoa(typeID) + "\n"
	}
	ioutil.WriteFile("typeDump.txt", []byte(f), 0777)
}

// PrintHeader prints a simple header about SDETool
func PrintHeader() {
	fmt.Println("SDETool written by Nick Powell")
	fmt.Println("The multiplatform CLI tool for accessing data from the Dust514 SDE")
}

// ForcePanic is a function to call panic() used for debuging our panic
// recovery and logging
func ForcePanic() {
	panic("Forced runtime panic")
}

// SetCWD sets our current working directory to our application directory.
// will make the folder if it doesn't already exist
func SetCWD() {
	err := os.Mkdir(os.Getenv("HOME")+"/.SDETool/", 0777)
	err1 := os.Chdir(os.Getenv("HOME") + "/.SDETool/")
	if err != nil {
		fmt.Println(err.Error())
	}
	if err1 != nil {
		fmt.Println(err.Error())
	}
}

// Uninstall deletes SDETool from our path if installed using methods that
// will be available in the future if anyone finds this useful :P
func Uninstall() {
	var fname string
	switch runtime.GOOS {
	case "Windows":
		fname = "SDETool.exe"
	default:
		fname = "SDETool"
	}
	err1 := os.Remove(fname)
	if err1 != nil {
		LErr(err1.Error())
		return
	}
	s, _ := exec.LookPath("SDETool")
	if s == "" {
		fmt.Println("No install of SDETool found\nUninstall aborted.")
		return
	}
	err := os.Remove(s)
	if err != nil {
		LErr(err.Error())
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Would you like to cleanup the database files, cache and config?\n[y/n]")
	var a string
	fmt.Scanf("%s", &a)
	if strings.ToLower(a) == "y" {
		fmt.Println("Deleting SDETool data")
		err1 := os.Remove(os.Getenv("HOME") + "/.SDETool/")
		if err1 != nil {
			fmt.Println("Unable to cleanup SDETool data:/")
			LErr(err1.Error())
			fmt.Println(err1.Error())
			return
		}
	} else {
		fmt.Println("DONE")
	}
}

// Inverse makes any bool the opposite
// True  | False
// False | True
func Inverse(b bool) bool {
	if b {
		return false
	} else {
		return true
	}
}
