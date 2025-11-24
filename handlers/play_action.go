package handlers

import (
	"net/http"
)

// Traite une action de jeu (ex: reset, coup joué)
func PlayActionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		action := r.FormValue("action")
		if action == "reset" {
			http.Redirect(w, r, "/game/init", http.StatusSeeOther)
			return
		}
		// Autres actions possibles
		http.Redirect(w, r, "/game/play", http.StatusSeeOther)
	}
}
