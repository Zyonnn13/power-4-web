package models

import "time"

// GameRecord stores a finished game's summary for the scoreboard/history.
type GameRecord struct {
	Player1 string    `json:"player1"`
	Player2 string    `json:"player2"`
	Winner  string    `json:"winner"` // empty string for draw
	Date    time.Time `json:"date"`
	Turns   int       `json:"turns"`
}

// History holds all finished games. It's a package-level variable so handlers can append to it.
// For this small app we don't enforce heavy concurrency controls; append is okay for now.
var History []GameRecord

// AddRecord appends a finished game record to the global history.
func AddRecord(r GameRecord) {
	History = append(History, r)
}
