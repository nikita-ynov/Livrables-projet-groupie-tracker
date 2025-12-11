package router

import (
	"net/http"
	"nft/controller"

	"github.com/gorilla/mux"
)

func InitRouter() *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("static"))))

	// Route Index
	r.HandleFunc("/", controller.Home)
	r.HandleFunc("/collections", controller.Collections)
	r.HandleFunc("/favorites", controller.Favorites)
	r.HandleFunc("/about", controller.About)
	r.HandleFunc("/item/{id}", controller.Id)

	return r
}
