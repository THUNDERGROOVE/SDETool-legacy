package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
)

var db *sql.DB

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

// Call to initialize the SDE
func DBInitialize() {
	var err error
	db, err = sql.Open("sqlite3", SDEFile)
	if err != nil {
		// Don't panic or exit, could just mean the SDE wasn't downloaded
		log.Error("Error opening the SDE ", err.Error())
		fmt.Println("Error opening the SDE", err.Error())
	}
}

// Returns an SDEType by typeID
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

func GetSDEWhereNameContains(name string) []SDEType {
	var typelist []SDEType
	rows, err := db.Query("SELECT * FROM CatmaAttributes WHERE catmaAttributeName == 'mDisplayName' AND catmaValueText LIKE '%" + name + "%'")
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

		typelist = append(typelist, GetSDETypeID(nTypeID))
	}
	return typelist
}

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
