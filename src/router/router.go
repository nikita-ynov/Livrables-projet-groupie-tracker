package router

import (
	"net/http"
	"nft/controller"
)

func InitRouter() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Route Index
	mux.HandleFunc("/", controller.Home)
	mux.HandleFunc("/collections", controller.Collections)
	mux.HandleFunc("/favorites", controller.Favorites)
	mux.HandleFunc("/about", controller.About)
	mux.HandleFunc("/:id", controller.Id)

	return mux
}
