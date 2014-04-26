/*
	Used to house all code related to grabbing fits from external location.
	This is just an example until G Torq adds proper sharing to his website.
*/
package extern

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Preset is a struct used in the parsing of CFL JSON
type Preset struct {
	Name    string `json:presetname`
	modules []Module
}

// Module is a struct used in the parsing of CFL JSON
type Module struct {
	TypeID   int    `json:typeid`
	SlotType string `json:slowtype`
	Index    int    `json:index`
}

// CFL is a struct used to Unmarshal CFL JSON from external locations.  This
// implementation is based primarily on G Torq's DUST CFL
type CFL struct {
	Version     string `json:clf-version`
	CLFType     string `json:X-clf-type`
	GeneratedBy string `json:X-generatedy`
	FittingID   int    `json:X-dust_thang_dk-ID`

	Metadata map[string]string `json:metadata`

	Ship    map[string]string `json:ship`
	Status  string            `json:status`
	Message string            `json:msg`
}

const (
	RESTBase = "http://dust.thang.dk/rest/getfittingCLF?FittingID="
)

// DTTKGetFittingsByID is a function to get a Dust.thang.tk shared fitting by
// it's fitting ID.
func DTTKGetFittingByID(id int) (CFL, error) {
	r, err := http.Get(RESTBase + strconv.Itoa(id))
	if err != nil {
		return CFL{}, err
	}
	b, _ := ioutil.ReadAll(r.Body)
	print(string(b), "\n")
	cfl := CFL{}
	err = json.Unmarshal(b, &cfl)
	return cfl, err
}
