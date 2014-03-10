/*
	math.go provides helpers for calculating various things
*/
package util

import (
	"math"
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
