/*
	The util package is for constants, package variables and functions that belong no where else.
	Used to consolidate most of the codebase out of the main package and allow other devs
	to use my code if they want.
*/
package util

import (
	"archive/zip"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	// Constants
	SDEFile = "dustSDE.db"                                                           // Name for the SDE database file to be used
	SDEUrl  = "http://cdn1.eveonline.com/community/DUST_SDE/Uprising_1.7_674383.zip" // URL to download the SDE
)

var (
	TimeFunc bool
)

// Used to measure execuation time of a function, use with defer
func TimeFunction(start time.Time, name string) {
	elapsed := time.Since(start)
	if TimeFunc {
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

// Used to make sure we have our SDEFile and if not we get it ourselves.
func CheckFile() {
	if _, err := os.Stat(SDEFile); os.IsNotExist(err) {
		fmt.Println("Downloading SDE (~1.8MB)")
		res, err := http.Get(SDEUrl)
		if err != nil {
			fmt.Println("Error downloading SDE")
			os.Exit(1)
		}
		cont, err2 := ioutil.ReadAll(res.Body)
		err3 := ioutil.WriteFile(SDEFile+".zip", cont, 0777)
		if err2 != nil || err3 != nil {
			fmt.Println("Error saving SDE zip to disk")
			os.Exit(1)
		}
		r, err4 := zip.OpenReader(SDEFile + ".zip")
		if err4 != nil {
			fmt.Println("Error unzipping SDE zip")
		}
		reader, err5 := r.File[0].Open()
		data, err6 := ioutil.ReadAll(reader)
		err7 := ioutil.WriteFile(SDEFile, data, 0777)
		if err5 != nil || err6 != nil || err7 != nil {
			fmt.Println("Error writing file ", err5.Error(), err6.Error(), err7.Error())
			os.Exit(1)
		}
		fmt.Println("Downloaded successfully")
		fmt.Println("Optimizing DB...")
		db, err9 := sql.Open("sqlite3", SDEFile)
		if err9 != nil {
			// Don't panic or exit, could just mean the SDE wasn't downloaded
			fmt.Println("Error opening the SDE", err9.Error())
		}
		defer db.Close()
		db.Exec("CREATE INDEX lookup ON CatmaAttributes(catmaAttributeName,catmaValueText)")
	} else {
		// Delete our SDEZip when we get a chance
		os.Remove(SDEFile + ".zip")
	}
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
