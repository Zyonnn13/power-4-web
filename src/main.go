package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"

	"power-4-web/models"
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
		q := r.URL.Query()
		p1 := q.Get("player1")
		p2 := q.Get("player2")
		winner := q.Get("winner")
		turns := 0
		if t := q.Get("turns"); t != "" {
			if v, err := strconv.Atoi(t); err == nil {
				turns = v
			}
		}

		rec := models.GameRecord{
			Player1: p1,
			Player2: p2,
			Winner:  winner,
			Date:    time.Now(),
			Turns:   turns,
		}
		models.AddRecord(rec)

		data := struct {
			Player1 string
			Player2 string
			Winner  string
			Date    string
			Turns   int
		}{
			Player1: rec.Player1,
			Player2: rec.Player2,
			Winner:  rec.Winner,
			Date:    rec.Date.Format("2006-01-02 15:04:05"),
			Turns:   rec.Turns,
		}

		temp.ExecuteTemplate(w, "end", data)
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
