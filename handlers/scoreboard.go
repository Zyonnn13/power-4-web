package handlers

import (
	"html/template"
	"net/http"
)

func ScoreboardHandler(temp *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Tu peux ajouter ici la logique pour charger les scores depuis models
		temp.ExecuteTemplate(w, "scoreboard", nil)
	}
}
