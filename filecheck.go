package main

import (
	"archive/zip"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"net/http"
	"os"
)

// Used to make sure we have our SDEFile and if not we get it ourselves.
func checkfile() {
	if _, err := os.Stat(SDEFile); os.IsNotExist(err) {
		fmt.Println("Downloading SDE (~1.8MB)")
		res, err := http.Get(SDEUrl)
		if err != nil {
			fmt.Println("Error downloading SDE")
			os.Exit(1)
		}
		cont, err2 := ioutil.ReadAll(res.Body)
		err3 := ioutil.WriteFile(SDEFile+".zip", cont, 0777)
		if err2 != nil || err3 != nil {
			fmt.Println("Error saving SDE zip to disk")
			os.Exit(1)
		}
		r, err4 := zip.OpenReader(SDEFile + ".zip")
		if err4 != nil {
			fmt.Println("Error unzipping SDE zip")
		}
		reader, err5 := r.File[0].Open()
		data, err6 := ioutil.ReadAll(reader)
		err7 := ioutil.WriteFile(SDEFile, data, 0777)
		if err5 != nil || err6 != nil || err7 != nil {
			fmt.Println("Error writing file ", err5.Error(), err6.Error(), err7.Error())
			os.Exit(1)
		}
		fmt.Println("Downloaded successfully")
		fmt.Println("Optimizing DB...")
		db, err9 := sql.Open("sqlite3", SDEFile)
		if err9 != nil {
			// Don't panic or exit, could just mean the SDE wasn't downloaded
			fmt.Println("Error opening the SDE", err9.Error())
		}
		defer db.Close()
		db.Exec("CREATE INDEX lookup ON CatmaAttributes(catmaAttributeName,catmaValueText)")
	}
}
