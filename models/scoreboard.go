package models

import "time"

type GameRecord struct {
	Winner string
	Loser  string
	Date   string
}

// Scoreboard stocke l'historique de toutes les parties
var Scoreboard []GameRecord

// AddRecord ajoute une partie terminée au tableau des scores
func AddRecord(game *Game) {
	winnerName := "Égalité"
	loserName := "Personne"

	if game.Status == "win" {
		winnerName = game.Winner

		if game.Players[0].Name == winnerName {
			loserName = game.Players[1].Name
		} else {
			loserName = game.Players[0].Name
		}
	} else {
		// Cas du match nul
		winnerName = game.Players[0].Name + " & " + game.Players[1].Name
	}

	record := GameRecord{
		Winner: winnerName,
		Loser:  loserName,
		Date:   time.Now().Format("02/01/2006 15:04"),
	}

	
	Scoreboard = append([]GameRecord{record}, Scoreboard...)
}