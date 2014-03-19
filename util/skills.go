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
		t.Attribs.CPU = parseFloat("mVICProp.maxCpuReserve", t)
		t.Attribs.PG = parseFloat("mVICProp.maxPowerReserve", t)
		t.Attribs.Armor = parseFloat("mVICProp.maxArmor", t)
		t.Attribs.Shields = parseFloat("mVICProp.maxShield", t)

		t.Attribs.ArmorRepair = parseFloat("mVICProp.healArmorRate", t)
		t.Attribs.ScanPrecision = parseFloat("mVICProp.signatureScanPrecision", t)
		t.Attribs.ScanProfile = parseFloat("mVICProp.signatureScanProfile", t)
		t.Attribs.ScanRadius = parseFloat("mVICProp.signatureScanRadius", t)
		t.Attribs.MeleeDamage = parseFloat("mCharMeleeProp.meleeDamage", t)
		t.Attribs.Stamina = parseFloat("mCharProp.maxStamina", t)
		t.Attribs.StaminaRecovery = parseFloat("mCharProp.staminaRecoveryPerSecond", t)
		t.Attribs.ShieldRechargeDelay = parseFloat("mVICProp.shieldRechargeDelay", t)
		t.Attribs.ShieldRechargeDepleted = parseFloat("mVICProp.shieldRechargePauseOnShieldDepleted", t)
		t.Attribs.ShieldRechargeRate = parseFloat("mVICProp.healShieldRate", t)
		t.Attribs.HackSpeedFactor = parseFloat("mHackSpeedFactor", t)
	} else if t.HasTag(category.Tag_weapon) {
		t.Attribs.AbsoluteRange = parseFloat("mFireMode0.absoluteRange", t)
		t.Attribs.EffectiveRange = parseFloat("mFireMode0.effectiveRange", t)
		t.Attribs.FireInterval = parseFloat("mFireMode0.fireInterval", t)
		t.Attribs.Damage = parseFloat("mFireMode0.instantHitDamage", t)
		t.Attribs.SplashDamage = parseFloat("mFireMode0.instantHitSplashDamage", t)
		t.Attribs.SplashRadius = parseFloat("mFireMode0.instantHitSplashDamageRadius", t)
		t.Attribs.CPU = parseFloat("mVICProp.amountCpuUsage", t)
		t.Attribs.PG = parseFloat("mVICProp.amountPowerUsage", t)
		t.Attribs.ShotCost, _ = strconv.Atoi(t.Attributes["mFireMode0.shotCost"])
		t.Attribs.ShotPerRound, _ = strconv.Atoi(t.Attributes["mFireMode0.shotPerRound"])
	}
}

func parseFloat(s string, t *SDEType) float64 {
	v, err := strconv.ParseFloat(t.Attributes[s], 64)
	if err != nil {
		LErr(fmt.Sprint("Error Parsing float from ", s, " value ", t.Attributes[s], " error: "+err.Error()))
	}
	return v
}

func (t *SDEType) applyAttributeToType(attribute string, value float64, method string, level int) {
	defer TimeFunction(time.Now(), t.GetName()+".applyAttributeToType("+attribute+", "+strconv.FormatFloat(value, 'f', 6, 64)+", "+method+", "+strconv.Itoa(level)+") "+method)
	Info("Applying attribute " + attribute)
	for k, _ := range t.Attributes {
		if k == attribute { // found
			ov := parseFloat(k, t)
			switch method {
			case "ADD":
				value = value * float64(level)
				value += float64(ov)
				Info(fmt.Sprint("Method ADD new value", value))
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
				if value > 1 {
					value = ov + (y * ov)
				} else {
					value = ov - (y * ov)
				}
				Info("Method: " + method + " Original value: " + strconv.FormatFloat(ov, 'f', 4, 64) + " v: " + strconv.FormatFloat(v, 'f', 4, 64) + " y: (v * float64(level)) / 100): " + strconv.FormatFloat(y, 'f', 4, 64) + " Output value: " + strconv.FormatFloat(value, 'f', 4, 64))
			default:
				LErr("Unknown attribute method" + method)
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
			value := parseFloat("modifier."+modint+".modifierValue", &b)
			if attrib == "" || method == "" {
				LErr("found broken modifer")
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
	t.ApplyAttributesToType()
}
