package util

/*
	skills.go is used for skill resolving and applying bonuses to types recursivly by finding skill dependencies.
	It will allow you to add a skill level to an SDEType which will automatically apply it to the type's values.
*/

import (
	"fmt"
	"github.com/THUNDERGROOVE/SDETool/category"
	"strconv"
	"strings"
	"time"
)

/*TODO
  1) Finish implementing none additive modifiers
    * This includes the modiferBonus(?) that some skills have
  2) Clean up our functions, they are very all over the place in naming and a lot of
     them don't seem to be needed more than once?
*/

// ApplyAttributesToType applies values from Attributes to Attrib to give
// ApplySkillsToType a safe place to modify values that can be accessed easily
func (t *SDEType) ApplyAttributesToType() { // t.Attribs. = t.Attributes[""]
	defer TimeFunction(time.Now(), t.GetName()+".ApplyAttributesToType()")
	if t.HasTag(category.Tag_dropsuit) {
		t.Attribs.CPU, _ = strconv.Atoi(t.Attributes["mVICProp.maxCpuReserve"])
		t.Attribs.PG, _ = strconv.Atoi(t.Attributes["mVICProp.maxPowerReserve"])
		t.Attribs.Armor, _ = strconv.Atoi(t.Attributes["mVICProp.maxArmor"])
		t.Attribs.Shields, _ = strconv.Atoi(t.Attributes["mVICProp.maxShields"])

		t.Attribs.ArmorRepair, _ = strconv.ParseFloat(t.Attributes["mVICProp.healArmorRate"], 64)
		t.Attribs.ScanPrecision, _ = strconv.ParseFloat(t.Attributes["mVICProp.signatureScanPrecision"], 64)
		t.Attribs.ScanProfile, _ = strconv.ParseFloat(t.Attributes["mVICProp.signatureScanProfile"], 64)
		t.Attribs.ScanRadius, _ = strconv.ParseFloat(t.Attributes["mVICProp.signatureScanRadius"], 64)
		t.Attribs.MeleeDamage, _ = strconv.ParseFloat(t.Attributes["mCharMeleeProp.meleeDamage"], 64)
		t.Attribs.Stamina, _ = strconv.ParseFloat(t.Attributes["mVICProf.maxStamina"], 64)
		t.Attribs.ShieldRechargeDelay, _ = strconv.ParseFloat(t.Attributes["mVICProp.shieldRechargeDelay"], 64)
		t.Attribs.ShieldRechargeDepleted, _ = strconv.ParseFloat(t.Attributes["mVICProp.shieldRechargePauseOnShieldDepleted"], 64)
		t.Attribs.ShieldRechargeRate, _ = strconv.ParseFloat(t.Attributes["mVICProp.healShieldRate"], 64)
		t.Attribs.HackSpeedFactor, _ = strconv.ParseFloat(t.Attributes["mHackSpeedFactor"], 64)
	} else if t.HasTag(category.Tag_weapon) {
		t.Attribs.CPU, _ = strconv.Atoi(t.Attributes["mVICProp.amountCpuUsage"])
		t.Attribs.PG, _ = strconv.Atoi(t.Attributes["mVICProp.amountPowerUsage"])
		t.Attribs.AbsoluteRange, _ = strconv.ParseFloat(t.Attributes["mFireMode0.absoluteRange"], 64)
		t.Attribs.EffectiveRange, _ = strconv.ParseFloat(t.Attributes["mFireMode0.effectiveRange"], 64)
		t.Attribs.FireInterval, _ = strconv.ParseFloat(t.Attributes["mFireMode0.fireInterval"], 64)
		t.Attribs.Damage, _ = strconv.ParseFloat(t.Attributes["mFireMode0.instantHitDamage"], 64)
		t.Attribs.SplashDamage, _ = strconv.ParseFloat(t.Attributes["mFireMode0.instantHitSplashDamage"], 64)
		t.Attribs.SplashRadius, _ = strconv.ParseFloat(t.Attributes["mFireMode0.instantHitSplashDamageRadius"], 64)
		t.Attribs.ShotCost, _ = strconv.Atoi(t.Attributes["mFireMode0.shotCost"])
		t.Attribs.ShotPerRound, _ = strconv.Atoi(t.Attributes["mFireMode0.shotPerRound"])
	}
}

func (t *SDEType) applyAttributeToType(attribute string, value float64, method string, level int) {
	defer TimeFunction(time.Now(), t.GetName()+".applyAttributeToType("+attribute+", "+strconv.FormatFloat(value, 'f', 6, 64)+", "+method+", "+strconv.Itoa(level)+")")
	fmt.Println("Applying attribute", attribute)
	for k, _ := range t.Attributes {
		if k == attribute { // found
			ov, _ := strconv.ParseFloat(t.Attributes[k], 64)
			switch method {
			case "ADD":
				value = value * float64(level)
				value += float64(ov)
			case "SUB":
				value = value * float64(level)
				value -= float64(ov)
			case "MULTIPLY":
				v := float64(0)
				if value > 1 {
					v = (100 * value) - 100
				} else {
					v = 100 - (100 * value)
				}
				y := (v * float64(level)) / 100
				value = ov - (y * ov)
			default:
				fmt.Println("Unknown attribute method", method)
			}
			t.Attributes[k] = strconv.FormatFloat(value, 'f', 2, 64)
		}
	}
	t.ApplyAttributesToType() // Fix values in t.Attribs
}

func (t *SDEType) getAllSkills() {
	defer TimeFunction(time.Now(), t.GetName()+".getAllSkills()")
	// Start by adding all required skills to the skills used for our type
	for k, _ := range t.Attributes {
		if strings.Contains(k, "requiredSkills") {
			skillint := strings.Split(k, ".")[1]
			tid, _ := strconv.Atoi(t.Attributes["requiredSkills."+skillint+".skillTypeID"])
			a := GetSDETypeIDFast(tid)
			a.getAllSkills()
			sname := a.GetName()
			if t.Skills[sname] == "" { // Assume already doesn't exist
				t.Skills[sname] = t.Attributes["requiredSkills."+skillint+".skillTypeID"] + "|" + t.Attributes["requiredSkills."+skillint+".skillLevel"] // allows for easy spliting  typenames won't have | in it.  May change to 0x0 but this works

			}

		}
	}
	for _, v := range t.Skills {
		stype := strings.Split(v, "|")[0]
		ltype := strings.Split(v, "|")[1]
		tid, _ := strconv.Atoi(stype)
		level, _ := strconv.Atoi(ltype)
		t.skillApply(tid, level) // apply all skill values

	}
}

func (t *SDEType) skillApply(skillTID int, level int) {
	defer TimeFunction(time.Now(), t.GetName()+".skillApply("+strconv.Itoa(skillTID)+")")
	b := GetSDETypeID(skillTID)
	// Apply all modifier.X
	for k, _ := range b.Attributes {
		if strings.Contains(k, ".modifierType") {
			modint := strings.Split(k, ".")[1]
			attrib := b.Attributes["modifier."+modint+".attributeName"]
			method := b.Attributes["modifier."+modint+".modifierType"]
			value, _ := strconv.ParseFloat(b.Attributes["modifier."+modint+".modifierValue"], 64)
			if attrib == "" || method == "" {
				fmt.Println("ERROR: Found broken modifer")
				continue
			}
			t.applyAttributeToType(attrib, value, method, level)
		}
	}
}

// ApplySkillsToType searches each type for the required skills and applies skill
// bonuses.  It also looks for skills required by each skill until we have none left.
func (t *SDEType) ApplySkillsToType() {
	defer TimeFunction(time.Now(), t.GetName()+".ApplySkillsToType()")
	Skills := make(map[string]string)

	for k, _ := range t.Attributes {
		if strings.Contains(k, "requiredSkills") {
			skillint := strings.Split(k, ".")[1]
			tid, _ := strconv.Atoi(t.Attributes["requiredSkills."+skillint+".skillTypeID"])
			a := GetSDETypeIDFast(tid)
			sname := a.GetName()
			Skills[sname] = t.Attributes["requiredSkills."+skillint+".skillTypeID"] + "|" + t.Attributes["requiredSkills."+skillint+".skillLevel"]
		}
	}
	t.Skills = Skills // add initial skills
	t.getAllSkills()  // finish getting rest of skills recusivlyi
	fmt.Println(t.Skills)
	t.ApplyAttributesToType()
}
