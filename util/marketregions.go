package util

/*
	market.regions.go is just a hardcoded anonymous struct to keep track of all of the regions
*/

type Location struct {
	Name   string
	TypeID int
}
type Locations struct {
	Regions []Location
}

var Regions Locations

func init() {
	reg := make(map[string]int)
	Regions = Locations{}
	reg["Derelik"] = 10000001
	reg["The Forge"] = 10000002
	reg["Vale of the Silent"] = 10000003
	reg["UUA-F4"] = 10000004
	reg["Detorid"] = 10000005
	reg["Wicked Creek"] = 10000006
	reg["Cache"] = 10000007
	reg["Scalding Pass"] = 10000008
	reg["Insmother"] = 10000009
	reg["Tribute"] = 10000010
	reg["Great Wildlands"] = 10000011
	reg["Curse"] = 10000012
	reg["Malpais"] = 10000013
	reg["Catch"] = 10000014
	reg["Venal"] = 10000015
	reg["Lonetrek"] = 10000016
	reg["J7HZ-F"] = 10000017
	reg["The Spire"] = 10000018
	reg["A821-A"] = 10000019
	reg["Tash-Murkon"] = 10000020
	reg["Outer Passage"] = 10000021
	reg["Stain"] = 10000022
	reg["Pure Blind"] = 10000023
	reg["Immensea"] = 10000025
	reg["Etherium Reach"] = 10000027
	reg["Molden Heath"] = 10000028
	reg["Geminate"] = 10000029
	reg["Heimatar"] = 10000030
	reg["Impass"] = 10000031
	reg["Sinq Laison"] = 10000032
	reg["The Citadel"] = 10000033
	reg["The Kalevala Expanse"] = 10000034
	reg["Deklein"] = 10000035
	reg["Devoid"] = 10000036
	reg["Everyshore"] = 10000037
	reg["The Bleak Lands"] = 10000038
	reg["Esoteria"] = 10000039
	reg["Oasa"] = 10000040
	reg["Syndicate"] = 10000041
	reg["Metropolis"] = 10000042
	reg["Domain"] = 10000043
	reg["Solitude"] = 10000044
	reg["Tenal"] = 10000045
	reg["Fade"] = 10000046
	reg["Providence"] = 10000047
	reg["Placid"] = 10000048
	reg["Khanid"] = 10000049
	reg["Querious"] = 10000050
	reg["Cloud Ring"] = 10000051
	reg["Kador"] = 10000052
	reg["Cobalt Edge"] = 10000053
	reg["Aridia"] = 10000054
	reg["Branch"] = 10000055
	reg["Feythabolis"] = 10000056
	reg["Outer Ring"] = 10000057
	reg["Fountain"] = 10000058
	reg["Paragon Soul"] = 10000059
	reg["Delve"] = 10000060
	reg["Tenerifis"] = 10000061
	reg["Omist"] = 10000062
	reg["Period Basis"] = 10000063
	reg["Essence"] = 10000064
	reg["Kor-Azor"] = 10000065
	reg["Perrigen Falls"] = 10000066
	reg["Genesis"] = 10000067
	reg["Verge Vendor"] = 10000068
	reg["Black Rise"] = 10000069
	reg["A-R00001"] = 11000001
	reg["A-R00002"] = 11000002
	reg["A-R00003"] = 11000003
	reg["B-R00004"] = 11000004
	reg["B-R00005"] = 11000005
	reg["B-R00006"] = 11000006
	reg["B-R00007"] = 11000007
	reg["B-R00008"] = 11000008
	reg["C-R00009"] = 11000009
	reg["C-R00010"] = 11000010
	reg["C-R00011"] = 11000011
	reg["C-R00012"] = 11000012
	reg["C-R00013"] = 11000013
	reg["C-R00014"] = 11000014
	reg["C-R00015"] = 11000015
	reg["D-R00016"] = 11000016
	reg["D-R00017"] = 11000017
	reg["D-R00018"] = 11000018
	reg["D-R00019"] = 11000019
	reg["D-R00020"] = 11000020
	reg["D-R00021"] = 11000021
	reg["D-R00022"] = 11000022
	reg["D-R00023"] = 11000023
	reg["E-R00024"] = 11000024
	reg["E-R00025"] = 11000025
	reg["E-R00026"] = 11000026
	reg["E-R00027"] = 11000027
	reg["E-R00028"] = 11000028
	reg["E-R00029"] = 11000029
	reg["F-R00030"] = 11000030
	for k, v := range reg {
		k := Location{k, v}
		Regions.Regions = append(Regions.Regions, k)
	}
}
