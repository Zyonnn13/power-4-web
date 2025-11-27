package models

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type GameRecord struct {
	Player1 string
	Player2 string
	Winner  string
	Turns   int
	Date    time.Time
}

var Scoreboard []GameRecord

const dbFile = "scoreboard.json"

func AddRecord(game *Game) {
	turns := 0
	for r := 0; r < Rows; r++ {
		for c := 0; c < Cols; c++ {
			if game.Grid[r][c] != "" {
				turns++
			}
		}
	}

	winnerName := ""
	if game.Status == "win" {
		winnerName = game.Winner
	}

	record := GameRecord{
		Player1: game.Players[0].Name,
		Player2: game.Players[1].Name,
		Winner:  winnerName,
		Turns:   turns,
		Date:    time.Now(),
	}

	Scoreboard = append([]GameRecord{record}, Scoreboard...)

	SaveScoreboard()
}

func SaveScoreboard() {
	data, err := json.MarshalIndent(Scoreboard, "", "  ")
	if err != nil {
		fmt.Println("Erreur lors de la sauvegarde du scoreboard:", err)
		return
	}
	_ = os.WriteFile(dbFile, data, 0644)
}

func LoadScoreboard() {
	data, err := os.ReadFile(dbFile)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		fmt.Println("Erreur lors du chargement du scoreboard:", err)
		return
	}
	_ = json.Unmarshal(data, &Scoreboard)
}
