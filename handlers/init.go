package handlers

import (
	"html/template"
	"net/http"
)

// Affiche la page d'initialisation
func InitPageHandler(temp *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "init", nil)
	}
}
