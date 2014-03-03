package main

import (
	"fmt"
	"math"
	"strconv"
)

// Methods for figureing out damages

// Returns the multiplier to apply to a value when stacking
// S(n) = 0.5^[((n-1) / 2.22292081) ^2]
func StackingMultiplier(n int) float64 {
	// Written out so it makes sense to me.
	a := n - 1
	b := float64(a) / float64(2.22292081)
	c := math.Pow(b, 2)
	d := math.Pow(0.5, c)
	return d
}

// Returns raw damages with prof and damage mods taken into account
// Doesn't work on swarms or anything that gets damage from it's projectile
func (t *SDEType) GetRawDamage(ProfLvl, ComplexModCount, EnhancedModCount, BasicModCount int) float64 {
	baseDamage, err := strconv.ParseFloat(t.Attributes["mFireMode0.instantHitDamage"], 64)
	if err != nil {
		fmt.Println("Error converting mFireMode0.instantHitDamage to int", err.Error())
		return 0
	}
	k := 1
	Modifier := float64(0)
	for i := 0; i < ComplexModCount; i++ {
		Modifier += StackingMultiplier(k) * 10
		k++

	}
	for i := 0; i < EnhancedModCount; i++ {
		Modifier += StackingMultiplier(k) * 5
		k++

	}
	for i := 0; i < BasicModCount; i++ {
		Modifier += StackingMultiplier(k) * 3
		k++

	}
	for i := 0; i < ProfLvl; i++ {
		Modifier += 2

	}
	addDam := float64(baseDamage) * float64((Modifier / 100))
	return baseDamage + addDam
}
