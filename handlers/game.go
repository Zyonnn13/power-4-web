package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"power-4-web/models"
)

var (
	currentGame  *models.Game
	mutex        sync.Mutex
	gameRecorded bool
)

type PlayPageData struct {
	Game     models.Game
	Columns  []int
	Error    string
	Message  string
	CanReset bool
}

// InitPageHandler affiche le formulaire d'initialisation d'une partie.
func InitPageHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}
		if err := tmpl.ExecuteTemplate(w, "init", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// InitProcessHandler traite les données du formulaire d'initialisation.
func InitProcessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/game/init", http.StatusSeeOther)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Redirect(w, r, "/game/init?error=formulaire+invalide", http.StatusSeeOther)
			return
		}
		p1 := strings.TrimSpace(r.FormValue("player1"))
		p2 := strings.TrimSpace(r.FormValue("player2"))
		choice := r.FormValue("colorChoice")
		p1Color, p2Color := resolveColors(choice)

		mutex.Lock()
		currentGame = models.NewGame()
		currentGame.ConfigurePlayers(p1, p2, p1Color, p2Color)
		gameRecorded = false
		mutex.Unlock()

		http.Redirect(w, r, "/game/play", http.StatusSeeOther)
	}
}

func InitGame(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	currentGame = models.NewGame()
	gameRecorded = false
	w.WriteHeader(http.StatusOK)
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
	if err != nil || col < 0 || col >= models.Cols {
		http.Error(w, "Invalid column", http.StatusBadRequest)
		return
	}

	success := currentGame.DropToken(col)
	if !success {
		http.Error(w, "Column full or game ended", http.StatusBadRequest)
		return
	}

	// Renvoie l’état du jeu en JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentGame)
}

// PlayPageHandler renvoie l'interface HTML du jeu.
func PlayPageHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock()
		if currentGame == nil {
			currentGame = models.NewGame()
		}
		gameSnapshot := *currentGame
		data := PlayPageData{
			Game:     gameSnapshot,
			Columns:  buildColumns(),
			Error:    r.URL.Query().Get("error"),
			Message:  r.URL.Query().Get("info"),
			CanReset: gameSnapshot.Finished,
		}
		mutex.Unlock()

		if err := tmpl.ExecuteTemplate(w, "play", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// PlayActionHandler gère les coups et redirections associées.
func PlayActionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/game/play", http.StatusSeeOther)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Redirect(w, r, "/game/play?error=formulaire+invalide", http.StatusSeeOther)
			return
		}

		mutex.Lock()
		if currentGame == nil {
			currentGame = models.NewGame()
		}

		redirectURL := "/game/play"
		params := url.Values{}

		switch r.FormValue("action") {
		case "reset":
			prevP1 := currentGame.Player1Name
			prevP2 := currentGame.Player2Name
			prevC1 := currentGame.Player1Color
			prevC2 := currentGame.Player2Color
			currentGame = models.NewGame()
			currentGame.ConfigurePlayers(prevP1, prevP2, prevC1, prevC2)
			gameRecorded = false
			params.Set("info", "Nouvelle partie relancée avec les mêmes joueurs.")
		default:
			colStr := r.FormValue("col")
			col, err := strconv.Atoi(colStr)
			if err != nil || col < 0 || col >= models.Cols {
				params.Set("error", "Colonne invalide.")
			} else if !currentGame.DropToken(col) {
				params.Set("error", "Colonne pleine ou partie terminée.")
			} else if currentGame.Finished {
				if currentGame.Winner == 0 {
					params.Set("info", "Égalité ! La grille est pleine.")
				} else {
					params.Set("info", currentGame.PlayerName(currentGame.Winner)+" remporte la partie.")
				}
				redirectURL = "/game/end"
			}
		}
		mutex.Unlock()

		if query := params.Encode(); query != "" {
			redirectURL = redirectURL + "?" + query
		}
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}

type EndPageData struct {
	Player1 string
	Player2 string
	Winner  string
	Date    string
	Turns   int
	Message string
}

// EndPageHandler gère l'affichage de la fin de partie et ajoute l'entrée au scoreboard.
func EndPageHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock()
		if currentGame == nil || !currentGame.Finished {
			mutex.Unlock()
			http.Redirect(w, r, "/game/init", http.StatusSeeOther)
			return
		}

		if !gameRecorded {
			winner := ""
			if currentGame.Winner != 0 {
				winner = currentGame.PlayerName(currentGame.Winner)
			}
			record := models.GameRecord{
				Player1: currentGame.Player1Name,
				Player2: currentGame.Player2Name,
				Winner:  winner,
				Date:    time.Now(),
				Turns:   currentGame.Moves,
			}
			models.AddRecord(record)
			gameRecorded = true
		}

		data := EndPageData{
			Player1: currentGame.Player1Name,
			Player2: currentGame.Player2Name,
			Turns:   currentGame.Moves,
			Date:    time.Now().Format("02/01/2006 15:04"),
			Message: r.URL.Query().Get("info"),
		}
		if currentGame.Winner != 0 {
			data.Winner = currentGame.PlayerName(currentGame.Winner)
		}
		mutex.Unlock()

		if err := tmpl.ExecuteTemplate(w, "end", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

type ScoreboardData struct {
	Records []models.GameRecord
}

// ScoreboardHandler affiche l'historique des parties.
func ScoreboardHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := ScoreboardData{Records: snapshotHistory()}
		if err := tmpl.ExecuteTemplate(w, "scoreboard", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func buildColumns() []int {
	result := make([]int, models.Cols)
	for i := 0; i < models.Cols; i++ {
		result[i] = i
	}
	return result
}

func resolveColors(choice string) (string, string) {
	switch choice {
	case "yellowred":
		return "yellow", "red"
	default:
		return "red", "yellow"
	}
}

func snapshotHistory() []models.GameRecord {
	mutex.Lock()
	defer mutex.Unlock()
	history := make([]models.GameRecord, len(models.History))
	copy(history, models.History)
	for i, j := 0, len(history)-1; i < j; i, j = i+1, j-1 {
		history[i], history[j] = history[j], history[i]
	}
	return history
}
