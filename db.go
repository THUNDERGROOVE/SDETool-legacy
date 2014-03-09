/*
	db.go contains all of the methods related to handling the DustSDE database.
*/

package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"strings"
	"time"
)

var (
	db *sql.DB
)

// SDEType is used to help manipulate and look up information about a particular type
// in the SDE
type SDEType struct {
	TypeID         int
	TypeName       string
	Class          string
	Attributes     map[string]string
	Tags           []int
	HighModules    int
	LowModules     int
	LightWeapons   int
	HeavyWeapons   int
	EquipmentSlots int
	GrenadeSlots   int
	Sidearms       int
}

type CatmaAttributeLookup struct {
	nTypeID            int
	catmaAttributeName string
	catmaValue         string
}

// DBInitialize is used to to initialize the SDE database file
func DBInitialize() {
	var err error
	db, err = sql.Open("sqlite3", SDEFile)
	if err != nil {
		// Don't panic or exit, could just mean the SDE wasn't downloaded
		fmt.Println("Error opening the SDE", err.Error())
	}
}

// GetSDETypeID returns an SDEType by typeID
// Currently bottlenecking our GetSDEWhereNameContains takes ~ 300 ms to execute
func GetSDETypeID(TID int) SDEType {
	// Get ID
	var err error
	rows, err := db.Query("SELECT * FROM CatmaTypes WHERE typeID == " + strconv.Itoa(TID))
	if err != nil {
		fmt.Println(err.Error())
		return SDEType{}
	}
	var typeID int
	var typeName string
	rows.Next()
	err = rows.Scan(&typeID, &typeName)
	if err != nil {
		fmt.Println("Error getting type by TypeID", err.Error())
		return SDEType{}
	}
	sde := SDEType{}
	sde.TypeID = typeID
	sde.TypeName = typeName
	sde.Attributes = make(map[string]string)

	// Get attributes
	rows, err = db.Query("SELECT * FROM CatmaAttributes WHERE typeID == " + strconv.Itoa(typeID))
	if err != nil {
		fmt.Println("SELECT * FROM CatmaAttributes WHERE typeID == "+string(typeID)+"\n", err.Error())
		return SDEType{}
	}
	for rows.Next() {
		var nTypeID int
		var catmaAttributeName string
		var catmaValueInt string
		var catmaValueReal string

		var catmaValueText string

		err := rows.Scan(&nTypeID, &catmaAttributeName, &catmaValueInt, &catmaValueReal, &catmaValueText)

		if err != nil {
			fmt.Println("Error parsing attribute\n\t", err.Error())
			continue
		}

		if catmaValueInt != "None" {
			sde.Attributes[catmaAttributeName] = string(catmaValueInt)
		}
		if catmaValueReal != "None" {
			sde.Attributes[catmaAttributeName] = string(catmaValueReal)
		}
		if catmaValueText != "None" {
			sde.Attributes[catmaAttributeName] = string(catmaValueText)
		}

	}
	sde._GetTags()
	sde._GetModules()
	return sde
}

// GetSDETypeIDFast returns an SDEType by typeID
// Doesn't get tag or module info use when you just need base info
func GetSDETypeIDFast(TID int) SDEType {
	// Get ID
	var err error
	rows, err := db.Query("SELECT * FROM CatmaTypes WHERE typeID == " + strconv.Itoa(TID))
	if err != nil {
		fmt.Println(err.Error())
		return SDEType{}
	}
	var typeID int
	var typeName string
	rows.Next()
	err = rows.Scan(&typeID, &typeName)
	if err != nil {
		fmt.Println("Error getting type by TypeID", err.Error())
		return SDEType{}
	}
	sde := SDEType{}
	sde.TypeID = typeID
	sde.TypeName = typeName
	sde.Attributes = make(map[string]string)

	// Get attributes
	rows, err = db.Query("SELECT * FROM CatmaAttributes WHERE typeID == " + strconv.Itoa(typeID) + " AND catmaAttributeName == 'mDisplayName'")
	if err != nil {
		fmt.Println("SELECT * FROM CatmaAttributes WHERE typeID == "+string(typeID)+"\n", err.Error())
		return SDEType{}
	}
	for rows.Next() {
		var nTypeID int
		var catmaAttributeName string
		var catmaValueInt string
		var catmaValueReal string

		var catmaValueText string

		err := rows.Scan(&nTypeID, &catmaAttributeName, &catmaValueInt, &catmaValueReal, &catmaValueText)

		if err != nil {
			fmt.Println("Error parsing attribute\n\t", err.Error())
			continue
		}

		if catmaValueInt != "None" {
			sde.Attributes[catmaAttributeName] = string(catmaValueInt)
		}
		if catmaValueReal != "None" {
			sde.Attributes[catmaAttributeName] = string(catmaValueReal)
		}
		if catmaValueText != "None" {
			sde.Attributes[catmaAttributeName] = string(catmaValueText)
		}

	}
	return sde
}

// SearchSDE returns a slice of SDETypes by using either GetSDEWhereNameContains
// or GetSDEWhereNameContainsFast depending on how many values we _think_ are going
// to be returned
func SearchSDE(name string) []SDEType {
	defer timeFunction(time.Now(), "SearchSDE("+name+")")
	rows, err := db.Query("SELECT * FROM CatmaAttributes WHERE catmaValueText LIKE '%" + name + "%' AND catmaAttributeName == 'mDisplayName'")
	if err != nil {
		fmt.Println(err.Error())
		return []SDEType{}
	}
	counter := 0
	for rows.Next() {
		counter++
	}
	if counter > 16 {
		return GetSDEWhereNameContainsFast(name)
	} else {
		return GetSDEWhereNameContains(name)
	}

}

// GetSDEWhereNameContains returns a slice of SDETypes whose mDisplayName contains name
func GetSDEWhereNameContains(name string) []SDEType {
	defer timeFunction(time.Now(), "GetSDEWhereNameContains("+name+")")
	var typelist []SDEType
	rows, err := db.Query("SELECT * FROM CatmaAttributes WHERE catmaValueText LIKE '%" + name + "%' AND catmaAttributeName == 'mDisplayName'")
	if err != nil {
		fmt.Println(err.Error())
		return []SDEType{}
	}
	for rows.Next() {
		var nTypeID int
		var catmaAttributeName string
		var catmaValueInt string
		var catmaValueReal string
		var catmaValueText string

		err := rows.Scan(&nTypeID, &catmaAttributeName, &catmaValueInt, &catmaValueReal, &catmaValueText)

		if err != nil {
			fmt.Println("Error parsing attribute\n\t", err.Error())
			continue
		}

		typelist = append(typelist, GetSDETypeIDFast(nTypeID))
	}
	return typelist
}

// This is meant to be a much faster version of GetSDEWhereNameContains where
// We cache all of our CatmaAttributes and check that before attempting to
// attempt a SELECT statement.  Meant to be used for searches where you only
// need a typeID and mDisplayName.  It's only faster when there are more than
// ~16 values returned, the larger slice returned the better it is to use this
// function.
func GetSDEWhereNameContainsFast(name string) []SDEType {
	defer timeFunction(time.Now(), "GetSDEWhereNameContainsFast("+name+")")
	var typelist []SDEType
	var lookup []CatmaAttributeLookup

	rows, err := db.Query("SELECT * FROM CatmaAttributes")
	if err != nil {
		fmt.Println(err.Error())
		return typelist
	}

	for rows.Next() {
		var nTypeID int
		var catmaAttributeName string
		var catmaValueInt string
		var catmaValueReal string
		var catmaValueText string

		var catmaValue string

		err := rows.Scan(&nTypeID, &catmaAttributeName, &catmaValueInt, &catmaValueReal, &catmaValueText)

		if err != nil {
			fmt.Println("Error parsing attribute\n\t", err.Error())
			continue
		}
		if catmaValueInt != "None" {
			catmaValue = catmaValueInt
		} else if catmaValueReal != "None" {
			catmaValue = catmaValueReal
		} else if catmaValueText != "None" {
			catmaValue = catmaValueText
		}
		lookup = append(lookup, CatmaAttributeLookup{nTypeID, catmaAttributeName, catmaValue})
	}

	for _, v := range lookup {
		if v.catmaAttributeName == "mDisplayName" && strings.Contains(strings.ToLower(v.catmaValue), strings.ToLower(name)) {
			temp := SDEType{}
			temp.TypeID = v.nTypeID
			temp.Attributes = make(map[string]string)
			temp.Attributes["mDisplayName"] = v.catmaValue
			typelist = append(typelist, temp)
		}
	}
	return typelist
}

// GetTypeName returns the name of a type when given a TypeID
func GetTypeName(typeID int) string {
	rows, err := db.Query("SELECT * FROM CatmaTypes WHERE typeID == " + strconv.Itoa(typeID))
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	var typeName string
	rows.Next()
	err = rows.Scan(&typeID, &typeName)
	if err != nil {
		fmt.Println("Error getting type by TypeID", err.Error())
		return ""
	}
	return typeName
}

// _GetTags is a helper function to apply the tags of a type to an SDEType
// Bottlenecking, 100ms execution time
func (t *SDEType) _GetTags() {
	rows, err := db.Query("SELECT * FROM CatmaAttributes WHERE typeID == " + strconv.Itoa(t.TypeID) + " AND catmaAttributeName LIKE 'tag.%'")
	if err != nil {
		fmt.Println("Error getting tags", err.Error())
		return
	}
	for rows.Next() {
		var nTypeID int
		var catmaAttributeName string
		var catmaValueInt string
		var catmaValueReal string

		var catmaValueText string

		err := rows.Scan(&nTypeID, &catmaAttributeName, &catmaValueInt, &catmaValueReal, &catmaValueText)
		if err != nil {
			fmt.Println("Error parsing attribute\n\t", err.Error())
			continue
		}
		s, _ := strconv.Atoi(catmaValueInt)
		t.Tags = append(t.Tags, s)
	}
}

// _GetModules is a helper function to add the module counts to an SDEType
// Bottlenecking, 100ms execution time
func (t *SDEType) _GetModules() {
	rows, err := db.Query("SELECT * FROM CatmaAttributes WHERE typeID == " + strconv.Itoa(t.TypeID) + " AND catmaAttributeName LIKE 'mModuleSlots.%'")
	if err != nil {
		fmt.Println("Error getting tags", err.Error())
		return
	}
	for rows.Next() {
		var nTypeID int
		var catmaAttributeName string
		var catmaValueInt string
		var catmaValueReal string

		var catmaValueText string

		err := rows.Scan(&nTypeID, &catmaAttributeName, &catmaValueInt, &catmaValueReal, &catmaValueText)
		if err != nil {
			fmt.Println("Error parsing attribute\n\t", err.Error())
			continue
		}
		switch catmaValueText {
		case "IL":
			t.LowModules++
		case "IH":
			t.HighModules++
		case "GM":
			t.GrenadeSlots++
		case "IE":
			t.EquipmentSlots++
		case "WP":
			t.LightWeapons++
		case "WS":
			t.Sidearms++
		}
	}
}
