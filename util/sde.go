package util

import (
	//"database/sql"
	"fmt"
	"github.com/THUNDERGROOVE/SDETool/args"
	"github.com/THUNDERGROOVE/SDETool/category"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strconv"
	"strings"
	"time"
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
	LargeTurrets   int
	SmallTurrets   int
}

// GetRawDamage returns raw damages with prof and damage mods taken into account
// Doesn't work on swarms or anything that gets damage from it's projectile __yet__
func (t *SDEType) GetRawDamage(ProfLvl, ComplexModCount, EnhancedModCount, BasicModCount int) float64 {
	// Slice of ints, used to pass all the damage modifier values to GenericCalculateValue
	m := []int{}
	for i := 0; i < ComplexModCount; i++ {
		m = append(m, 10)
	}
	for i := 0; i < EnhancedModCount; i++ {
		m = append(m, 5)
	}
	for i := 0; i < BasicModCount; i++ {
		m = append(m, 3)
	}
	v, err := t.GenericCalculateValue("mFireMode0.instantHitDamage", true, []int{ProfLvl * 3}, m)
	if err != nil {
		fmt.Println(err.Error())
		return float64(-1) // Obvious error, don't see the need to have this method return an error at this time
	}
	return v
}

// GenericCalculateValue returns a float64 of a generic value calculated from various variables
// You can pass a slice of int slices of values to calculate, each slice
// within that slices' values will do stacking penalties
// HighOrLow should be set, true = high, false = low
// Requires that you give it at least a string that gets a value from t.Attributes
// which must be convertable to an int, will throw an error if unable to.
func (t *SDEType) GenericCalculateValue(ValueAttribute string, HighOrLow bool, Modifiers ...[]int) (float64, error) {
	BaseValue, err := strconv.ParseFloat(t.Attributes[ValueAttribute], 64)
	if err != nil {
		return BaseValue, err
	}
	Modifier := float64(0)

	for _, v := range Modifiers {
		Modifier = float64(0)
		for k, c := range v {
			// Apply our value
			if HighOrLow {
				Modifier += StackingMultiplier(k) * float64(c)
			} else {
				Modifier -= StackingMultiplier(k) * float64(c)
			}
		}
		BaseValue = BaseValue + (float64(BaseValue) * float64((Modifier / 100)))
	}
	return BaseValue, nil
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
	if t.HasTag(category.TagVehicle) {
		fmt.Println("===== Vehicle =====")
		fmt.Println("-> High slots:", t.HighModules)
		fmt.Println("-> Low slots:", t.LowModules)
		fmt.Println("-> Large Turrets:", t.LargeTurrets)
		fmt.Println("-> Small Turrets:", t.SmallTurrets)
	}
	if len(t.Tags) > 0 {

		fmt.Println("===== Tags =====")
		for _, c := range t.Tags {
			fmt.Println("->", c, GetTypeName(c))
		}
	}
	if *args.VerboseInfo == true {
		l := longestLen(t.Attributes)
		fmt.Println(l)
		if len(t.Attributes) > 0 {
			fmt.Println("===== Attributes =====:")
			for k, v := range t.Attributes {
				// Don't print descriptions
				if k == "mDescription" || k == "mShortDescription" {
					continue
				}
				fmt.Println(k + xspaces(longestLen(t.Attributes)-len(k)) + " | " + v)
			}
		} else {
			fmt.Println("No attributes to show")
		}
	}
}

// GetSDETypeID returns an SDEType by typeID
// Currently bottlenecking our GetSDEWhereNameContains takes ~ 300 ms to execute
func GetSDETypeID(TID int) SDEType {
	if TID == 0 {
		fmt.Println("Unable to find type by ID type lookup failed, make sure the name provided actually exists")
		os.Exit(1)
	}
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
		fmt.Println("Error getting type by TypeID, GetSDETypeID("+strconv.Itoa(TID)+")", err.Error())
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
		fmt.Println("Error getting type by TypeID, GetSDETypeIDFast("+strconv.Itoa(TID)+")", err.Error())
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
	defer TimeFunction(time.Now(), "SearchSDE("+name+")")
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
		return GetSDEWhereNameContainsFastC(name)
	} else {
		return GetSDEWhereNameContains(name)
	}

}

// GetSDEWhereNameContains returns a slice of SDETypes whose mDisplayName contains name
func GetSDEWhereNameContains(name string) []SDEType {
	defer TimeFunction(time.Now(), "GetSDEWhereNameContains("+name+")")
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
	defer TimeFunction(time.Now(), "GetSDEWhereNameContainsFast("+name+")")
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

// This is meant to be a much faster version of GetSDEWhereNameContainsFast
// however we use GoRoutines and channels to _hopefully_ speed up proccessing
// of our query.  Saves ~0.2 seconds on SDETool -s rail :D
func GetSDEWhereNameContainsFastC(name string) []SDEType {
	defer TimeFunction(time.Now(), "GetSDEWhereNameContainsFastC("+name+")")
	var typelist []SDEType
	var lookup []CatmaAttributeLookup
	SDETypeChan := make(chan SDEType)
	breakr := false

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

	go func() {
		for {
			select {
			case k := <-SDETypeChan:
				typelist = append(typelist, k)
			case <-time.Tick(100 * time.Millisecond):
				if breakr {
					break
				}
			}
		}
	}()

	for _, v := range lookup {
		if v.catmaAttributeName == "mDisplayName" && strings.Contains(strings.ToLower(v.catmaValue), strings.ToLower(name)) {
			temp := SDEType{}
			temp.TypeID = v.nTypeID
			temp.Attributes = make(map[string]string)
			temp.Attributes["mDisplayName"] = v.catmaValue
			SDETypeChan <- temp
		}
	}
	breakr = true
	return typelist
}

// GetTypeName returns the name of a type when given a TypeID
func GetTypeName(typeID int) string {
	defer TimeFunction(time.Now(), "GetTypeName("+strconv.Itoa(typeID)+")")
	rows, err := db.Query("SELECT * FROM CatmaTypes WHERE typeID == " + strconv.Itoa(typeID))
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	var typeName string
	rows.Next()
	err = rows.Scan(&typeID, &typeName)
	if err != nil {
		fmt.Println("Error getting type name from TypeID, GetTypeName("+strconv.Itoa(typeID)+")", err.Error())
		return ""
	}
	return typeName
}

// GetTypeIDByName returns the TypeID of a type when given a typeName
func GetTypeIDByName(typeName string) int {
	defer TimeFunction(time.Now(), "GetTypeIDByName("+typeName+")")
	rows, err := db.Query("SELECT * FROM CatmaTypes WHERE typeName == '" + typeName + "'")
	if err != nil {
		if strings.Contains(err.Error(), "no such column") {
			fmt.Println("Unable to find a type with the name", typeName)
			os.Exit(1)
		}
		return 0
	}
	var typeID int
	rows.Next()
	err = rows.Scan(&typeID, &typeName)
	if err != nil {
		fmt.Println("Error getting type by TypeID, GetTypeIDByName("+typeName+")", err.Error())
		return 0
	}
	return typeID
}

func GetTypeIDByDName(name string) int {
	defer TimeFunction(time.Now(), "GetTypeIDByDName("+name+")")
	pmatches := make(map[int]string)
	rows, err := db.Query("SELECT * FROM CatmaAttributes WHERE catmaAttributeName == 'mDisplayName' AND catmaValueText LIKE '%" + name + "%'")
	if err != nil {
		fmt.Println("Error getting type by display name", err.Error())
		return 0
	}
	for rows.Next() {

		var nTypeID int
		var catmaAttributeName string
		var catmaValueInt string
		var catmaValueReal string

		var catmaValueText string

		err = rows.Scan(&nTypeID, &catmaAttributeName, &catmaValueInt, &catmaValueReal, &catmaValueText)
		if err != nil {
			fmt.Println("Error parsing attribute\n\t", err.Error())
			return 0
		}
		pmatches[nTypeID] = catmaValueText
	}
	k := FuzzySearch(pmatches, name)
	return k
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
		case "VH":
			t.HighModules++
		case "VL":
			t.LowModules++
		case "GM":
			t.GrenadeSlots++
		case "IE":
			t.EquipmentSlots++
		case "WP":
			t.LightWeapons++
		case "WS":
			t.Sidearms++
		case "TL":
			t.LargeTurrets++
		case "TS":
			t.SmallTurrets++
		}
	}
	// Remove hidden slots
	for k, _ := range t.Attributes {
		if strings.Contains(k, "mModuleSlots") && strings.Contains(k, "slotType") {
			curSlot := strings.Split(k, ".")[1] // Hope no index issues
			if t.Attributes["mModuleSlots."+curSlot+".visible"] == "False" {
				switch t.Attributes["mModuleSlots."+curSlot+".slotType"] {
				case "IL":
					t.LowModules--
				case "IH":
					t.HighModules--
				case "VH":
					t.HighModules--
				case "VL":
					t.LowModules--
				case "GM":
					t.GrenadeSlots--
				case "IE":
					t.EquipmentSlots--
				case "WP":
					t.LightWeapons--
				case "WS":
					t.Sidearms--
				case "TL":
					t.LargeTurrets--
				case "TS":
					t.SmallTurrets--
				}
			}
		}
	}
}