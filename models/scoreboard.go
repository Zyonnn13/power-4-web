package models

import "time"

type GameRecord struct {
	Winner string
	Loser  string
	Date   string
}

var Scoreboard []GameRecord

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
		winnerName = game.Players[0].Name + " & " + game.Players[1].Name
	}

	record := GameRecord{
		Winner: winnerName,
		Loser:  loserName,
		Date:   time.Now().Format("02/01/2006 15:04"),
	}

	Scoreboard = append([]GameRecord{record}, Scoreboard...)
}
