package handlers

import (
	"net/http"
	"power-4-web/models"
)

func InitProcessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()

		player1 := r.FormValue("player1")
		player2 := r.FormValue("player2")

		currentGame = models.NewGame(player1, player2)

		http.Redirect(w, r, "/game/play", http.StatusSeeOther)
	}
}
