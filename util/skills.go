package util

/*
	skills.go is used for skill resolving and applying bonuses to types recursivly by finding skill dependencies.
	It will allow you to add a skill level to an SDEType which will automatically apply it to the type's values.
*/

import (
	"github.com/THUNDERGROOVE/SDETool/category"
	"strconv"
	"strings"
	"time"
)

/* TODO, finish the Skill type.
type Skill struct {
}*/

// ApplyAttributesToType applies values from Attributes to Attrib to give
// ApplySkillsToType a safe place to modify values that can be accessed easily
func (t *SDEType) ApplyAttributesToType() { // t.Attribs. = t.Attributes[""]
	defer TimeFunction(time.Now(), t.GetName()+".ApplyAttributesToType()")
	if t.HasTag(category.TagDropsuit) {
		t.Attribs.CPU, _ = strconv.Atoi(t.Attributes["mVICProp.maxCpuReserve"])
		t.Attribs.PG, _ = strconv.Atoi(t.Attributes["mVICProp.maxPowerReserve"])
		t.Attribs.ArmorRepair, _ = strconv.Atoi(t.Attributes["mVICProp.healArmorRate"])
		t.Attribs.Armor, _ = strconv.Atoi(t.Attributes["mVICProp.maxArmor"])
		t.Attribs.Shields, _ = strconv.Atoi(t.Attributes["mVICProp.maxShields"])
		t.Attribs.ScanPrecision, _ = strconv.Atoi(t.Attributes["mVICProp.signatureScanPrecision"])
		t.Attribs.ScanProfile, _ = strconv.Atoi(t.Attributes["mVICProp.signatureScanProfile"])
		t.Attribs.ScanRadius, _ = strconv.Atoi(t.Attributes["mVICProp.signatureScanRadius"])
		t.Attribs.MeleeDamage, _ = strconv.Atoi(t.Attributes["mCharMeleeProp.meleeDamage"])
		t.Attribs.Stamina, _ = strconv.Atoi(t.Attributes["mVICProf.maxStamina"])
		t.Attribs.ShieldRechargeDelay, _ = strconv.Atoi(t.Attributes["mVICProp.shieldRechargeDelay"])
		t.Attribs.ShieldRechargeDepleted, _ = strconv.Atoi(t.Attributes["mVICProp.shieldRechargePauseOnShieldDepleted"])
		t.Attribs.ShieldRechargeRate, _ = strconv.Atoi(t.Attributes["mVICProp.healShieldRate"])
		t.Attribs.HackSpeedFactor, _ = strconv.Atoi(t.Attributes["mHackSpeedFactor"])
	} else if t.HasTag(category.Tag_weapon) {
		t.Attribs.CPU, _ = strconv.Atoi(t.Attributes["mVICProp.amountCpuUsage"])
		t.Attribs.PG, _ = strconv.Atoi(t.Attributes["mVICProp.amountPowerUsage"])
		t.Attribs.AbsoluteRange, _ = strconv.Atoi(t.Attributes["mFireMode0.absoluteRange"])
		t.Attribs.EffectiveRange, _ = strconv.Atoi(t.Attributes["mFireMode0.effectiveRange"])
		t.Attribs.FireInterval, _ = strconv.Atoi(t.Attributes["mFireMode0.fireInterval"])
		t.Attribs.Damage, _ = strconv.Atoi(t.Attributes["mFireMode0.instantHitDamage"])
		t.Attribs.SplashDamage, _ = strconv.Atoi(t.Attributes["mFireMode0.instantHitSplashDamage"])
		t.Attribs.SplashRadius, _ = strconv.Atoi(t.Attributes["mFireMode0.instantHitSplashDamageRadius"])
		t.Attribs.ShotCost, _ = strconv.Atoi(t.Attributes["mFireMode0.shotCost"])
		t.Attribs.ShotPerRound, _ = strconv.Atoi(t.Attributes["mFireMode0.shotPerRound"])
	}
}

func (t *SDEType) applyAttributeToType(attribute string, value int, method string, level int) {
	defer TimeFunction(time.Now(), t.GetName()+".applyAttributeToType("+attribute+", "+strconv.Itoa(value)+", "+method+", "+strconv.Itoa(level)+")")
	for k, _ := range t.Attributes {
		if k == attribute { // found
			ov, _ := strconv.Atoi(t.Attributes[k])
			value = value * level
			if method == "ADD" {
				value += ov
			} else {
				value -= ov
			}
			t.Attributes[k] = strconv.Itoa(value)
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
			value, _ := strconv.Atoi(b.Attributes["modifier."+modint+".modifierValue"])
			if attrib == "" || method == "" || value == 0 {
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
	t.getAllSkills()  // finish getting rest of skills recusivly
}
