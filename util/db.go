package util

/*
	db.go contains all of the methods related to handling the DustSDE database.
*/
import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

const (
	SDEFile = "dustSDE.db"
)

var (
	db *sql.DB
)

type CatmaAttributeLookup struct {
	nTypeID            int
	catmaAttributeName string
	catmaValue         string
}

// DBInitialize is used to to initialize the SDE database file
func DBInitialize() {
	var err error
	db, err = sql.Open("sqlite3", SDEFile)
	if err != nil {
		// Don't panic or exit, could just mean the SDE wasn't downloaded
		fmt.Println("Error opening the SDE", err.Error())
	}
}
