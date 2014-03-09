package main

import (
	"github.com/THUNDERGROOVE/SDETool/category"
	"math"
	"testing"
)

func roundStackingMultipler(x float64) int {
	return int(math.Floor(x * 100))
}

func TestStackingMultiplier(t *testing.T) {
	for i := 0; i < 5; i++ { // Only potential for 5 damage mods
		switch i {
		case 1:
			if roundStackingMultipler(StackingMultiplier(i)) != 100 {
				t.Fatal("Stacking multipler failed for", i, "got", roundStackingMultipler(StackingMultiplier(i)))
				t.FailNow()
			}
		case 2:
			if roundStackingMultipler(StackingMultiplier(i)) != 86 {
				t.Fatal("Stacking multipler failed for", i, "got", roundStackingMultipler(StackingMultiplier(i)))
				t.FailNow()
			}
		case 3:
			if roundStackingMultipler(StackingMultiplier(i)) != 57 {
				t.Fatal("Stacking multipler failed for", i, "got", roundStackingMultipler(StackingMultiplier(i)))
				t.FailNow()
			}
		case 4:
			if roundStackingMultipler(StackingMultiplier(i)) != 28 {
				t.Fatal("Stacking multipler failed for", i, "got", roundStackingMultipler(StackingMultiplier(i)))
				t.FailNow()
			}
		case 5:
			if roundStackingMultipler(StackingMultiplier(i)) != 10 {
				t.Fatal("Stacking multipler failed for", i, "got", roundStackingMultipler(StackingMultiplier(i)))
				t.FailNow()
			}
		}
	}
}
