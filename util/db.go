package util

/*
	db.go contains all of the methods related to handling the DustSDE database.
*/

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
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
		fmt.Println("Error opening the SDE", err.Error())
		LErr("failed to open SDE DB")
		os.Exit(1)
	}
}
