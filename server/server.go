/*
	The server package encapsolates all of the code related to hosting the web verison of SDETool
*/
package server

import (
	"fmt"
	"github.com/THUNDERGROOVE/SDETool/args"
	"github.com/THUNDERGROOVE/SDETool/util"
	"html/template"
	"net/http"
	"strconv"
)

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		util.Info(fmt.Sprintf("%s %s %s", r.RemoteAddr, r.Method, r.URL))
		handler.ServeHTTP(w, r)
	})
}

func home(w http.ResponseWriter, r *http.Request) {
	homet, err := template.ParseFiles("server/templates/home.html")
	if err != nil {
		util.LErr(err.Error())
	}
	homet.Execute(w, nil)
}
func process(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Placeholder"))
}

func RunServer() {
	http.HandleFunc("/", home)
	http.HandleFunc("/process", process)
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "server"+r.URL.Path)
	})
	http.ListenAndServe(":"+strconv.Itoa(*args.Port), Log(http.DefaultServeMux))
}
