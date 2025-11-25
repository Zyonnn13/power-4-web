package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"power-4-web/handlers"
)

func main() {

	temp := template.Must(template.New("").Funcs(template.FuncMap{
		"seq": func(count int) []int {
			s := make([]int, count)
			for i := range s {
				s[i] = i
			}
			return s
		},
		"inc": func(i int) int {
			return i + 1
		},
	}).ParseGlob("./templates/*.html"))

	initPageHandler := handlers.InitPageHandler(temp)
	http.HandleFunc("/game/init", initPageHandler)
	http.HandleFunc("/init", initPageHandler)
	http.HandleFunc("/game/init/traitement", handlers.InitProcessHandler())

	playPageHandler := handlers.PlayPageHandler(temp)
	http.HandleFunc("/game/play", playPageHandler)
	http.HandleFunc("/play", playPageHandler)
	http.HandleFunc("/game/play/traitement", handlers.PlayActionHandler())

	endHandler := handlers.EndPageHandler(temp)
	http.HandleFunc("/game/end", endHandler)
	http.HandleFunc("/templates/end", endHandler)

	scoreboardHandler := handlers.ScoreboardHandler(temp)
	http.HandleFunc("/game/scoreboard", scoreboardHandler)
	http.HandleFunc("/templates/scoreboard", scoreboardHandler)

	http.HandleFunc("/api/init", handlers.InitGameAPI)
	http.HandleFunc("/api/play", handlers.PlayMove)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() != "/" {
			w.WriteHeader(http.StatusNotFound)
			temp.ExecuteTemplate(w, "404", nil)
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
