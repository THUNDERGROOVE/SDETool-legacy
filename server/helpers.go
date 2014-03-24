package server

import (
	// "github.com/THUNDERGROOVE/SDETool/args"
	"archive/zip"
	"bytes"
	"database/sql"
	"fmt"
	"github.com/THUNDERGROOVE/SDETool/util"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func init() {
	checkfile()
	util.DBInitialize()
}

var (
	doneOut = false
)

const (
	SDEFile = "dustSDE.db"
	SDEUrl  = "http://cdn1.eveonline.com/community/DUST_SDE/Uprising_1.7_674383.zip"
)

type function func()

func lbToBr(s string) string {
	return strings.Replace(strings.Replace(s, "\n", "<br>", -1), " ", "&nbsp;&nbsp;", -1)
}
func routeout(f function) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()
	w.Close()
	os.Stdout = old
	out := <-outC
	util.Info(out)
	return out
}
func doneroute() {
	doneOut = true
}

// Used to make sure we have our SDEFile and if not we get it ourselves.
func checkfile() {
	if _, err := os.Stat(SDEFile); os.IsNotExist(err) {
		util.TimeFunc = true
		defer util.TimeFunction(time.Now(), "checkfile() download")
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
	} else {
		// Delete our SDEZip when we get a chance
		os.Remove(SDEFile + ".zip")
	}
}

func processArgString(Args string) map[string]string {
	f := []byte(Args)
	t := make(map[string]string)
	for i := 0; i < len(Args); i++ {
		if f[i] == '-' {
			setIsVal := false
			setname := ""
			argname := ""
			for j := i + 1; j < len(f); j++ {
				if f[j] != ' ' && setIsVal == false {
					argname += string(f[j])
				} else if setIsVal == true && f[j] != '-' {
					setname += string(f[j])
				} else if f[j] == '-' {
					break
				} else {
					setIsVal = true
				}
			}
			t[argname] = setname
		}
	}
	return t
}
