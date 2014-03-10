/*
	util.go is for constants, package variables and functions that belong no where else
*/

package util

import (
	"fmt"
	"github.com/THUNDERGROOVE/SDETool/args"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Used to measure execuation time of a function, use with defer
func TimeFunction(start time.Time, name string) {
	elapsed := time.Since(start)
	if *args.TimeExecution {
		fmt.Printf("%s took %s\n", name, elapsed)
	}
}

const (
	sqlliteDriver = `SDETool uses an SQLite3 driver to deal with the SDE database from https://github.com/mattn/go-sqlite3, see http://mattn.mit-license.org/2012 for licensing`
)

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
