package controller

import (
	"net/http"
	"nft/pages"
)

func renderPage(w http.ResponseWriter, filename string, data any) {
	err := pages.Temp.ExecuteTemplate(w, filename, data)
	if err != nil {
		http.Error(w, "Erreur rendu template : "+err.Error(), http.StatusInternalServerError)
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "index.html", nil)
}

func Collections(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "collections.html", nil)
}

func Id(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "id.html", nil)
}

func Favorites(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "favorites.html", nil)
}

func About(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "about.html", nil)
}
