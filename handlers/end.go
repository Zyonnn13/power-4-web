package handlers

import (
	"html/template"
	"net/http"
)

// Affiche la page de fin de partie
func EndPageHandler(temp *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "end", nil)
	}
}
