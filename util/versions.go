package util

import (
	"archive/zip"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var (
	SDEVersion string
	SDEFile    = "dustSDE.db" // Name for the SDE database file to be used
	SDEUrl     = SDE1_8       // URL to download the SDE
)

const (
	// Constants
	SDE1_7 = "http://cdn1.eveonline.com/community/DUST_SDE/Uprising_1.7_674383.zip"
	SDE1_8 = "http://cdn1.eveonline.com/community/DUST_SDE/Uprising_1.8_739147.zip"
)

func init() {
	SDEFile = "dustSDE.db" // Name for the SDE database file to be used
	SDEUrl = SDE1_8        // URL to download the SDE
}

// Used to make sure we have our SDEFile and if not we get it ourselves.
func CheckFile() {
	switch SDEVersion {
	case "1.7":
		SDEUrl = SDE1_7
		SDEFile = strings.Join(strings.Split(SDEFile, "."), SDEVersion+".")
	case "1.8":
		SDEUrl = SDE1_8
		SDEFile = strings.Join(strings.Split(SDEFile, "."), SDEVersion+".")
	}
	if _, err := os.Stat(SDEFile); os.IsNotExist(err) {
		fmt.Println("Downloading SDE (~1.8MB)")
		res, err := http.Get(SDEUrl)
		if err != nil {
			LErr("downloading SDE")
			os.Exit(1)
		}
		cont, err2 := ioutil.ReadAll(res.Body)
		err3 := ioutil.WriteFile(SDEFile+".zip", cont, 0777)
		if err2 != nil || err3 != nil {
			LErr("saving SDE zip to disk")
			os.Exit(1)
		}
		r, err4 := zip.OpenReader(SDEFile + ".zip")
		if err4 != nil {
			LErr("unzipping SDE zip")
		}
		reader, err5 := r.File[0].Open()
		data, err6 := ioutil.ReadAll(reader)
		err7 := ioutil.WriteFile(SDEFile, data, 0777)
		if err5 != nil || err6 != nil || err7 != nil {
			LErr("writing db file")
			os.Exit(1)
		}
		fmt.Println("Downloaded successfully")
		fmt.Println("Optimizing DB...")
		db, err9 := sql.Open("sqlite3", SDEFile)
		if err9 != nil {
			// Don't panic or exit, could just mean the SDE wasn't downloaded
			LErr("calling sql.open on " + SDEFile)
		}
		defer db.Close()
		db.Exec("CREATE INDEX lookup ON CatmaAttributes(catmaAttributeName,catmaValueText)")
	} else {
		// Delete our SDEZip when we get a chance
		os.Remove(SDEFile + ".zip")
	}
}
