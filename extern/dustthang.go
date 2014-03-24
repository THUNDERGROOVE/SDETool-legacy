package extern

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Preset struct {
	Name    string `json:presetname`
	modules []Module
}
type Module struct {
	TypeID   int    `json:typeid`
	SlotType string `json:slowtype`
	Index    int    `json:index`
}
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

func DTDKGetFittingByID(id int) (CFL, error) {
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
