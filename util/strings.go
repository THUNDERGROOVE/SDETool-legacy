/*
	strings.go implements various functions for working with strings that aren't in the strings package
*/
package util

import (
	"strings"
)

// Returns the "closest" matching string to subs from strings as a TypeID
func FuzzySearch(Strings map[int]string, subs string) int {
	bm := 0
	bmMatchType := 0 // 0 none, 1 = strings.Contains() match, 2 non case sensitive exact match, 3 case sensitive exact match
	for k, v := range Strings {
		if strings.Contains(v, subs) && bmMatchType < 1 {
			bm = k
			bmMatchType = 1
		} else if strings.ToLower(v) == strings.ToLower(subs) && bmMatchType < 2 {
			bm = k
			bmMatchType = 2
		} else if v == subs && bmMatchType < 3 {
			bm = k
			bmMatchType = 3
		}
	}
	return bm
}
