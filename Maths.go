package main

import (
	"fmt"
	"math"
	"strconv"
)

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
