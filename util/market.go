package util

/*
	market.go handles getting data from  http://public_crest_sisi.testeveonline.com/market/
*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	//"os"
	"strconv"
	"time"
)

// MarketData is a set to group a slice of MarketDataEntry
type MarketData struct {
	Items []MarketDataEntry `json:"items"`
}

// MarketDataEntry is a struct for Unmarhaling Market data
type MarketDataEntry struct {
	AveragePrice int    `json:"avgPrice"`
	Date         string `json:"date"`
	HighPrice    int    `json:"highPrice"`
	LowPrice     int    `json:"lowPrice"`
	OrderCount   int    `json:"orderCount"`
	Volume       int    `json:"volume"`
	//OrderCountString string `json:"orderCount_str"`
	//VolumeString     string `json:"volume_str"`
}

// GetTotalISKSpent is a function to get the total amount of ISK spent on a certain type.
// It's currently broken until I look at market data again
func (s *SDEType) GetTotalISKSpent() int {
	defer TimeFunction(time.Now(), "(s *SDEType) GetTotalISKSpent()")
	fmt.Println("GetTotalISKSpent() called\n\n\n")
	t := s.TypeID
	TotalVolume := 0
	fmt.Println(len(Regions.Regions), "regions")
	for _, l := range Regions.Regions {
		v := l.TypeID
		fmt.Println("Region, " + strconv.Itoa(v))
		r, err := http.Get("http://public-crest.eveonline.com/market/" + strconv.Itoa(v) + "/types/" + strconv.Itoa(t) + "/history/")
		if err != nil || r == nil {
			fmt.Println("Error getting http://public-crest.eveonline.com/market/" + strconv.Itoa(v) + "/types/" + strconv.Itoa(t) + "/history/")
			continue
			//os.Exit(1)
		}
		a, err2 := ioutil.ReadAll(r.Body)
		if err2 != nil || a == nil {
			fmt.Println("Error reading from r.Body")
			continue
		}
		var Data MarketData
		err3 := json.Unmarshal(a, &Data)
		if err3 != nil {
			fmt.Println("Error unmarshaling data,", err3.Error())
			fmt.Println("Dumping errorounes JSON")
			ioutil.WriteFile("MarketJSONError.json", a, 0777)
			continue
		}
		fmt.Println(Data)
		for _, v := range Data.Items {
			fmt.Println(v.Date, "has", v.Volume, "items")
			TotalVolume += v.Volume
		}
	}
	p, _ := strconv.Atoi(s.GetPrice())
	return TotalVolume * p
}
