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

func checkfile() {
	if _, err := os.Stat(SDEFile); os.IsNotExist(err) {
		fmt.Println("Downloading SDE (~1.8MB)")
		res, err := http.Get(SDEUrl)
		if err != nil {
			// Log print and exit with error code
			fmt.Println("Error downloading SDE")
			log.Critical("Error downloading SDE")
			os.Exit(1)
		}
		cont, err2 := ioutil.ReadAll(res.Body)
		err3 := ioutil.WriteFile(SDEFile+".zip", cont, 0777)
		if err2 != nil || err3 != nil {
			// Log print and exit with error code
			fmt.Println("Error saving SDE zip to disk")
			log.Critical("Error saving SDE zip to disk")
			os.Exit(1)
		}
		r, err4 := zip.OpenReader(SDEFile + ".zip")
		if err4 != nil {
			log.Critical("Error unzipping SDE zip")
		}
		reader, err5 := r.File[0].Open()
		data, err6 := ioutil.ReadAll(reader)
		err7 := ioutil.WriteFile(SDEFile, data, 0777)
		if err5 != nil || err6 != nil || err7 != nil {
			// Log print and exit with error code
			fmt.Println("Error writing file ", err5.Error(), err6.Error(), err7.Error())
			log.Critical("Error writing file", err5.Error(), err6.Error(), err7.Error())
			os.Exit(1)
		}
		fmt.Println("Downloaded successfully")
		log.Info("Downloaded SDE successfully")
		err8 := os.Remove(SDEFile + ".zip")
		if err8 != nil {
			// Log print and exit with error code
			fmt.Println("Error Deleting zip", err8.Error())
			log.Critical("Error Deleting zip", err8.Error())
			os.Exit(1)
		}
		fmt.Println("Optimizing DB...")
		db, err9 := sql.Open("sqlite3", SDEFile)
		if err9 != nil {
			// Don't panic or exit, could just mean the SDE wasn't downloaded
			log.Error("Error opening the SDE ", err9.Error())
			fmt.Println("Error opening the SDE", err9.Error())
		}
		defer db.Close()
		db.Exec("CREATE INDEX lookup ON CatmaAttributes(catmaAttributeName,catmaValueText)")
	}
}
