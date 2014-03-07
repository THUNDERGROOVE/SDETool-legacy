package main

import (
	"fmt"
	"github.com/THUNDERGROOVE/SDETool/category"
	"strconv"
)

// Returns a string of spaces with the amount x
func xspaces(x int) string {
	var t string
	for i := 0; i < x; i++ {
		t = t + " "
	}
	return t
}

func (t *SDEType) HasTag(tag int) bool {
	for _, c := range t.Tags {
		if tag == c {
			return true
		}
	}
	return false
}

// Returns mDisplayName
func (t *SDEType) GetName() string {
	if t.Attributes["mDisplayName"] == "" {
		return GetTypeName(t.TypeID)
	}
	return t.Attributes["mDisplayName"]
}

// Returns mDescription
func (t *SDEType) GetDescription() string {
	return t.Attributes["mDescription"]
}

// Returns Short Description
func (t *SDEType) GetShortDescription() string {
	return t.Attributes["mShortDescription"]
}

// Returns price
func (t *SDEType) GetPrice() string {
	return t.Attributes["basePrice"]
}

// Is consumable?
func (t *SDEType) IsConsumable() bool {
	if t.Attributes["consumable"] == "True" {
		return true
	}
	return false
}

// Returns a types Catagory TypeID
func (t *SDEType) Category() int {
	c, err := strconv.Atoi(t.Attributes["categoryID"])
	if err != nil {
		return -1
	}
	return c
}

// Generic function to print the info of an SDEType
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
