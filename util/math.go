package util

/*
	math.go provides helpers for calculating various things
*/

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

// Helpers

// valToPercent converts 1.05 to 5 or 0.95 to 95
func valToPercent(v int) int {
	if v > 1 {
		return (v - 1) * 100
	} else {
		return v * 100
	}
}

// Methods for figureing out damages

// StackingMultiplier returns the multiplier to apply to a value when stacking
// S(n) = 0.5^[((n-1) / 2.22292081) ^2]
func StackingMultiplier(n int) float64 {
	// Written out so it makes sense to me.
	a := n - 1
	b := float64(a) / float64(2.22292081)
	c := math.Pow(b, 2)
	d := math.Pow(0.5, c)
	return d
}

// GetRawDamage returns raw damages with prof and damage mods taken into account
// Doesn't work on swarms or anything that gets damage from it's projectile __yet__
func (t *SDEType) GetRawDamage(ProfLvl, ComplexModCount, EnhancedModCount, BasicModCount int) float64 {
	fmt.Println(t.Attributes["mFireMode0.instantHitDamage"])
	t.ApplySkillsToType()
	t.ApplyAttributesToType()
	defer TimeFunction(time.Now(), fmt.Sprint("GetRawDamage(", ProfLvl, ComplexModCount, EnhancedModCount, BasicModCount, ")"))
	for i := 0; i < ComplexModCount; i++ {
		t.ApplyModuleToSuit(GetSDETypeID(351681)) // Complex mod typeID
	}
	for i := 0; i < EnhancedModCount; i++ {
		t.ApplyModuleToSuit(GetSDETypeID(351680)) // Enhanced mod typeID
	}
	for i := 0; i < BasicModCount; i++ {
		t.ApplyModuleToSuit(GetSDETypeID(351679)) // Basic mod typeID
	}
	return t.Attribs.Damage
}

// GenericCalculateValue returns a float64 of a generic value calculated from various variables
// You can pass a slice of int slices of values to calculate, each slice
// within that slices' values will do stacking penalties
// HighOrLow should be set, true = high, false = low
// Requires that you give it at least a string that gets a value from t.Attributes
// which must be convertable to an int, will throw an error if unable to.
func (t *SDEType) GenericCalculateValue(ValueAttribute string, HighOrLow bool, Modifiers ...[]int) (float64, error) {
	defer TimeFunction(time.Now(), fmt.Sprint("GenericCalculateValue(", ValueAttribute, HighOrLow, Modifiers, ")"))
	LErr("SDEType.GenericCalculateValue() is depreceated.  Consider using SDEType.ApplyModuleToSuit() instead")
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
		Info("GenericCalculateValue() called with a modifer of " + strconv.FormatFloat(float64((Modifier)), 'f', 3, 64))
		BaseValue = BaseValue + (float64(BaseValue) * float64((Modifier / 100)))
	}
	return BaseValue, nil
}
