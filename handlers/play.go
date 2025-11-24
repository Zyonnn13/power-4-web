package handlers

import (
	"html/template"
	"net/http"
)

// Affiche la page de jeu
func PlayPageHandler(temp *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "play", nil)
	}
}
