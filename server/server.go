/*
	The server package encapsolates all of the code related to hosting the web verison of SDETool
*/
package server

import (
	"fmt"
	"github.com/THUNDERGROOVE/SDETool/args"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hey there, this will soon be the home page of an SDETool server.  However I need to finish it :P")
}

func RunServer() {
	http.HandleFunc("/", home)
	http.ListenAndServe(":"+strconv.Itoa(*args.Port), nil)
}
