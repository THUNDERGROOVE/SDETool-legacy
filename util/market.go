package util

/*
	market.go handles getting data from  http://public_crest_sisi.testeveonline.com/market/
*/

import (
	"encoding/json"
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
	AveragePrice float64 `json:"avgPrice"`
	Date         string  `json:"date"`
	HighPrice    float64 `json:"highPrice"`
	LowPrice     float64 `json:"lowPrice"`
	OrderCount   float64 `json:"orderCount"`
	Volume       float64 `json:"volume"`
	//OrderCountString string `json:"orderCount_str"`
	//VolumeString     string `json:"volume_str"`
}

// GetTotalISKSpent is a function to get the total amount of ISK spent on a certain type.
// It's currently broken until I look at market data again
func (s *SDEType) GetTotalISKSpent() int {
	defer TimeFunction(time.Now(), "(s *SDEType) GetTotalISKSpent()")
	t := s.TypeID
	TotalVolume := float64(0)
	Info(len(Regions.Regions), "regions")
	for _, l := range Regions.Regions {
		v := l.TypeID
		r, err := http.Get("http://public-crest.eveonline.com/market/" + strconv.Itoa(v) + "/types/" + strconv.Itoa(t) + "/history/")
		if err != nil || r == nil {
			LErr("Error getting http://public-crest.eveonline.com/market/" + strconv.Itoa(v) + "/types/" + strconv.Itoa(t) + "/history/")
			continue
			//os.Exit(1)
		}
		a, err2 := ioutil.ReadAll(r.Body)
		if err2 != nil || a == nil {
			LErr("Error reading from r.Body")
			continue
		}
		var Data MarketData
		err3 := json.Unmarshal(a, &Data)
		if err3 != nil {
			LErr("Error unmarshaling data,", err3.Error())
			LErr("Dumping errorounes JSON")
			ioutil.WriteFile("MarketJSONError.json", a, 0777)
			continue
		}
		for _, v := range Data.Items {
			//Info(v.Date, "has", v.Volume, "items")
			TotalVolume += v.Volume
		}
		if len(Data.Items) == 0 {
			Info("Region: " + l.Name + " has no buy orders, this could mean this region doesn't sell DUST items")
		}
	}
	p, _ := strconv.Atoi(s.GetPrice())
	return int(TotalVolume * float64(p))
}
