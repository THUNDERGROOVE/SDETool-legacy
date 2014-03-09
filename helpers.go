package main

import (
	"fmt"
	"github.com/THUNDERGROOVE/SDETool/category"
	"strconv"
	"time"
)

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

// HasTag returns true if an SDEType contains a tag by TypeID
func (t *SDEType) HasTag(tag int) bool {
	for _, c := range t.Tags {
		if tag == c {
			return true
		}
	}
	return false
}

// GetName returns mDisplayName
func (t *SDEType) GetName() string {
	if t.Attributes["mDisplayName"] == "" {
		return GetTypeName(t.TypeID)
	}
	return t.Attributes["mDisplayName"]
}

// GetDescription returns mDescription
func (t *SDEType) GetDescription() string {
	return t.Attributes["mDescription"]
}

// GetShortDescription returns Short Description
func (t *SDEType) GetShortDescription() string {
	return t.Attributes["mShortDescription"]
}

// GetPrice returns price
func (t *SDEType) GetPrice() string {
	return t.Attributes["basePrice"]
}

// IsConsumable returns if a SDEType is consumable
func (t *SDEType) IsConsumable() bool {
	if t.Attributes["consumable"] == "True" {
		return true
	}
	return false
}

// Category returns a types Category TypeID
func (t *SDEType) Category() int {
	c, err := strconv.Atoi(t.Attributes["categoryID"])
	if err != nil {
		return -1
	}
	return c
}

// PrintInfo is a generic function to print the info of an SDEType
func (t *SDEType) PrintInfo() {
	fmt.Println("Getting stats on " + t.GetName())
	if t.GetDescription() != "" {
		fmt.Println("===== Description =====")
		fmt.Println(t.GetDescription())
	}
	if t.GetPrice() != "" {
		fmt.Println("-> Cost", t.GetPrice(), "ISK")
	}
	//  Scanner
	if t.Category() == category.ActiveScanner {
		fmt.Println("====== Scanner ======")
		fmt.Println("-> Scan DB", t.Attributes["activeScanSignaturePrecision"])
	}
	if t.HasTag(category.TagDropsuit) {
		fmt.Println("===== Dropsuit =====")
		fmt.Println("-> Heavy Weapons:", t.HeavyWeapons)
		fmt.Println("-> Light Weapons:", t.LightWeapons)
		fmt.Println("-> Sidearms:", t.Sidearms)
		fmt.Println("-> Equipment slots:", t.EquipmentSlots)
		fmt.Println("-> High slots:", t.HighModules)
		fmt.Println("-> Low slots:", t.LowModules)

	}
	if len(t.Tags) > 0 {

		fmt.Println("===== Tags =====")
		for _, c := range t.Tags {
			fmt.Println("->", c, GetTypeName(c))
		}
	}
	if *VerboseInfo == true {
		if len(t.Attributes) > 0 {
			fmt.Println("===== Attributes =====:")
			for k, v := range t.Attributes {
				fmt.Println(k + " | " + v)
			}
		} else {
			fmt.Println("No attributes to show")
		}
	}
}

// PrintLicense is a simple function to just print our licensing for everything
func PrintLicense() {
	fmt.Println("SDETool is under the MIT license.  See LICENSE for more info")
	fmt.Println(sqlliteDriver)
}

// Used to measure execuation time of a function, use with defer
func timeFunction(start time.Time, name string) {
	elapsed := time.Since(start)
	if *TimeExecution {
		fmt.Printf("%s took %s\n", name, elapsed)
	}
}