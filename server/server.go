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
	flags := r.FormValue("data")
	if flags == "" {
		w.Write([]byte("This is for robots, not you human scum!"))
		return
	}
	k := ""
	a := processArgString(flags)
	if a["s"] != "" {
		k = lbToBr(util.SearchSDEFlag(a["s"]))
	} else if a["i"] != "" {
		tid := util.ResolveInput(a["i"])
		t := util.GetSDETypeID(tid)
		k = lbToBr(t.StringInfo())
	} else if a["d"] != "" {
		w.Write(
			[]byte(
				lbToBr(
					routeout(
						func() {
							t := util.GetSDETypeID(util.ResolveInput(a["d"]))
							w.Write([]byte("Getting damage on: " + t.GetName() + "\n"))
							c, _ := strconv.Atoi(a["c"])
							e, _ := strconv.Atoi(a["e"])
							b, _ := strconv.Atoi(a["b"])
							p, _ := strconv.Atoi(a["p"])
							if c == 0 && e == 0 && b == 0 {
								t.PrintDamageChart()
								return
							}
							d := t.GetRawDamage(p, c, e, b)
							fmt.Println("->", t.GetName(), "would do ", d, "damage")
							if *args.Prof != 0 {
								fmt.Println("->", "Proficiency level", p)
							}
							if c != 0 {
								fmt.Println("->", c, "Complex damage modifiers")
							}
							if e != 0 {
								fmt.Println("->", e, "Complex damage modifiers")
							}
							if b != 0 {
								fmt.Println("->", c, "Complex damage modifiers")
							}
						}))))
	}
	w.Write([]byte(k))
}

func RunServer() {
	http.HandleFunc("/", home)
	http.HandleFunc("/process", process)
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "server"+r.URL.Path)
	})
	http.ListenAndServe(":"+strconv.Itoa(*args.Port), Log(http.DefaultServeMux))
}
