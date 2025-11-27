package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"power-4-web/models"
	"strconv"
	"sync"
	"time"
)

var (
	currentGame *models.Game
	mutex       sync.Mutex
)

func InitPageHandler(temp *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		temp.ExecuteTemplate(w, "init", nil)
	}
}

func StartGameLogic(p1, p2 string, colorChoice string) {
	mutex.Lock()
	defer mutex.Unlock()

	if p1 == "" {
		p1 = "Joueur 1"
	}
	if p2 == "" {
		p2 = "Joueur 2"
	}

	var p1Color, p2Color string

	if colorChoice == "yellowred" {
		p1Color = "yellow"
		p2Color = "red"
	} else {

		p1Color = "red"
		p2Color = "yellow"
	}

	currentGame = models.NewGame()

	currentGame.ConfigurePlayers(p1, p2, p1Color, p2Color)
}

func InitProcessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		p1 := r.FormValue("player1")
		p2 := r.FormValue("player2")

		colorChoice := r.FormValue("colorChoice")

		StartGameLogic(p1, p2, colorChoice)

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

		currentGame.ErrorMsg = ""

		colStr := r.FormValue("col")
		if colStr == "" {
			colStr = r.URL.Query().Get("col")
		}

		col, err := strconv.Atoi(colStr)
		if err == nil && col >= 0 && col < models.Cols {
			succes := currentGame.DropToken(col)
			if !succes {
				currentGame.ErrorMsg = "Colonne pleine"
			}
		}

		http.Redirect(w, r, "/game/play", http.StatusSeeOther)
	}
}

func EndPageHandler(temp *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock()
		defer mutex.Unlock()

		if currentGame == nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Calcul du nombre de tours joués (en comptant les pions dans la grille)
		turnsPlayed := 0
		for r := 0; r < models.Rows; r++ {
			for c := 0; c < models.Cols; c++ {
				if currentGame.Grid[r][c] != "" {
					turnsPlayed++
				}
			}
		}

		// Préparation des données pour le HTML
		data := struct {
			Winner  string
			Player1 string
			Player2 string
			Date    string
			Turns   int
			Message string
		}{
			Player1: currentGame.Players[0].Name,
			Player2: currentGame.Players[1].Name,
			Date:    time.Now().Format("02/01/2006 15:04"),
			Turns:   turnsPlayed,
		}

		if currentGame.Status == "win" {
			data.Winner = currentGame.Winner
			data.Message = "Félicitations !"
		} else if currentGame.Status == "draw" {
			data.Message = "Match nul, bien joué à tous les deux !"
		}

		temp.ExecuteTemplate(w, "end", data)
	}
}

func RestartHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock()
		defer mutex.Unlock()

		if currentGame != nil {

			p1Name := currentGame.Players[0].Name
			p2Name := currentGame.Players[1].Name

			colorChoice := "redyellow"
			if currentGame.Players[0].Color == "yellow" {
				colorChoice = "yellowred"
			}

			mutex.Unlock()
			StartGameLogic(p1Name, p2Name, colorChoice)
			mutex.Lock()
		}

		http.Redirect(w, r, "/game/play", http.StatusSeeOther)
	}
}

func ScoreboardHandler(temp *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		data := struct {
			Records []models.GameRecord
		}{
			Records: models.Scoreboard,
		}
		temp.ExecuteTemplate(w, "scoreboard", data)
	}
}

func InitGameAPI(w http.ResponseWriter, r *http.Request) {
	p1 := r.URL.Query().Get("player1")
	p2 := r.URL.Query().Get("player2")

	colorChoice := r.URL.Query().Get("colorChoice")

	StartGameLogic(p1, p2, colorChoice)

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
