package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"power-4-web/models"
	"strconv"
	"sync"
)

var (
	currentGame *models.Game
	mutex       sync.Mutex
)


func StartGameLogic(p1, p2 string) {
	mutex.Lock()
	defer mutex.Unlock()

	if p1 == "" {
		p1 = "Joueur 1"
	}
	if p2 == "" {
		p2 = "Joueur 2"
	}

	currentGame = models.NewGame()
	currentGame.ConfigurePlayers(p1, p2, "red", "yellow")
}

func InitPageHandler(temp *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "init", nil)
	}
}

func InitProcessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		p1 := r.FormValue("player1")
		p2 := r.FormValue("player2")
		StartGameLogic(p1, p2)
		http.Redirect(w, r, "/game/play", http.StatusSeeOther)
	}
}

func PlayPageHandler(temp *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock()
		defer mutex.Unlock()

		if currentGame == nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		
		if currentGame.Status != "playing" {
			http.Redirect(w, r, "/game/end", http.StatusSeeOther)
			return
		}
		

		cols := make([]int, models.Cols)
		for i := range cols {
			cols[i] = i
		}

		data := struct {
			Game    *models.Game
			Columns []int
		}{
			Game:    currentGame,
			Columns: cols,
		}

		temp.ExecuteTemplate(w, "play", data)
	}
}

func PlayActionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock()
		defer mutex.Unlock()

		if currentGame == nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		
		action := r.FormValue("action")
		if action == "reset" {
			p1Name := currentGame.Players[0].Name
			p2Name := currentGame.Players[1].Name

			
			currentGame = models.NewGame()
			currentGame.ConfigurePlayers(p1Name, p2Name, "red", "yellow")

			
			http.Redirect(w, r, "/game/play", http.StatusSeeOther)
			return
		}
	

		colStr := r.FormValue("col")
		if colStr == "" {
			colStr = r.URL.Query().Get("col")
		}

		col, err := strconv.Atoi(colStr)
		if err == nil && col >= 0 && col < models.Cols {
			currentGame.DropToken(col)
		}

		http.Redirect(w, r, "/game/play", http.StatusSeeOther)
	}
}

func EndPageHandler(temp *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock()
		defer mutex.Unlock()

		temp.ExecuteTemplate(w, "end", currentGame)
	}
}

func ScoreboardHandler(temp *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Assurez-vous que models.Scoreboard est bien défini et public dans votre package models
		temp.ExecuteTemplate(w, "scoreboard", models.Scoreboard)
	}
}

func InitGameAPI(w http.ResponseWriter, r *http.Request) {
	p1 := r.URL.Query().Get("player1")
	p2 := r.URL.Query().Get("player2")
	StartGameLogic(p1, p2)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Game Initialized"))
}

func PlayMove(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	if currentGame == nil {
		http.Error(w, "Game not initialized", http.StatusBadRequest)
		return
	}
	colStr := r.URL.Query().Get("col")
	col, err := strconv.Atoi(colStr)
	if err != nil {
		http.Error(w, "Invalid column", http.StatusBadRequest)
		return
	}
	success := currentGame.DropToken(col)
	if !success {
		http.Error(w, "Column full or game ended", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentGame)
}