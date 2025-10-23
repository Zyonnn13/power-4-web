package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

func main() {

	temp, errtemp := template.ParseGlob("./templates/*.html")
	if errtemp != nil {
		fmt.Println(errtemp)
		os.Exit(1)
	}

	http.HandleFunc("/templates/init", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "init", nil)
	})

	http.HandleFunc("/templates/play", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "play", nil)
	})

	http.HandleFunc("/templates/end", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "end", nil)
	})

	http.HandleFunc("/templates/scoreboard", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "scoreboard", nil)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() != "/" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404"))
			return
		}
		temp.ExecuteTemplate(w, "index", nil)
	})

	chemin, _ := os.Getwd()
	fmt.Println(chemin)
	fileserver := http.FileServer(http.Dir(chemin + "/assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fileserver))

	http.ListenAndServe(":8000", nil)
}
