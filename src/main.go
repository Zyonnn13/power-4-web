package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

var temp *template.Template

func main() {

	// Charger tous les templates HTML
	var err error
	temp = template.New("")
	temp, err = temp.ParseGlob("./templates/*.html")
	if err != nil {
		fmt.Println("Erreur chargement templates:", err)
		os.Exit(1)
	}
	temp, err = temp.ParseGlob("./templates/game/*.html")
	if err != nil {
		fmt.Println("Erreur chargement templates game:", err)
		os.Exit(1)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./assets"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() != "/" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404"))
			return
		}
		temp.ExecuteTemplate(w, "index", nil)
	})

	fmt.Println("Serveur démarré sur http://localhost:8000")
	http.ListenAndServe(":8000", nil)
}
