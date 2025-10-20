package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

var temp = template.Must(template.ParseGlob("templates/*"))

func main() {

	temp, errtemp := template.ParseGlob("./assets/temp/*.html")
	if errtemp != nil {
		fmt.Println(errtemp)
		os.Exit(1)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() != "/" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404"))
			return
		}
		temp.ExecuteTemplate(w, "index", products)
	})

	http.ListenAndServe(":8080", nil)
}
