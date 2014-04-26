package util

import (
	"github.com/THUNDERGROOVE/SDETool/category"
	"strconv"
	"strings"
	"time"
)

/*
	The modules file houses methods for applying modules to dropsuits which takes into account skills
*/

// ApplyModuleToSuit applied a module by typeID to a suit.  This can be used
// for other uses as type application works very simiarly for most things.
func (t *SDEType) ApplyModuleToSuit(m SDEType) {
	defer TimeFunction(time.Now(), t.GetName()+".ApplyModuleToSuit("+m.GetName()+")")
	if t.HasTag(category.Tag_dropsuit) == false && t.HasTag(category.Tag_weapon) == false { // Allow to work with weapons
		LErr("attempted to call ApplyModuleToSuit() to non dropsuit type")
		return
	}
	if m.HasTag(category.Tag_module) == false {
		LErr("attempted to call ApplyModuleToSuit() with a non module type, looking for " + strconv.Itoa(category.Tag_module))
		for _, v := range t.Tags {
			LErr("Found tag " + strconv.Itoa(v))
		}
		return
	}
	m.ApplySkillsToType() // Apply skill bonuses
	t.moduleApply(m)
	t.Modules = append(t.Modules, m)

}

// moduleApply is a helper function.  It expects skills be applied beforehand
func (t *SDEType) moduleApply(b SDEType) {
	defer TimeFunction(time.Now(), t.GetName()+".moduleApply("+b.GetName()+")")
	// Apply all modifier.X
	for k, _ := range b.Attributes {
		if strings.Contains(k, ".modifierType") {
			modint := strings.Split(k, ".")[1]
			attrib := b.Attributes["modifier."+modint+".attributeName"]
			method := b.Attributes["modifier."+modint+".modifierType"]
			value, _ := strconv.ParseFloat(b.Attributes["modifier."+modint+".modifierValue"], 64)
			if attrib == "" || method == "" {
				LErr("found broken modifer")
				continue
			}
			ModStackCount := 0
			for _, v := range t.Modifiers {
				if attrib == v {
					ModStackCount++
				}
			}
			ModStackCount++
			t.applyAttributeToType(attrib, value, method, 1, ModStackCount)
		}
	}
}

// modulesAreValid checks if an SDEType with modules applied would be valid.
// For example a militia Gallente medium frame can't use a basic shield
// extender as it has no high slots
func (t *SDEType) ModulesAreValid() bool {
	defer TimeFunction(time.Now(), t.GetName()+".moduleApply()")
	IL := 0
	IH := 0
	VH := 0
	VL := 0
	for _, v := range t.Modules {
		switch v.Attributes["slotType"] {
		case "IL":
			IL++
		case "IH":
			IH++
		case "VH":
			VH++
		case "VL":
			VL++
		default:
			LErr("unknown slotType " + v.Attributes["slotType"])
		}
	}
	if VH > t.HighModules || VL > t.LowModules || IH > t.HighModules || IL > t.LowModules {
		return false
	}
	return true
}
