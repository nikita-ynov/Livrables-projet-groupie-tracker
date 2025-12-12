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

	// --- ROUTES PAGES (HTML) ---
	r.HandleFunc("/", controller.Home).Methods("GET")
	r.HandleFunc("/collections", controller.Collections).Methods("GET")
	r.HandleFunc("/favorites", controller.Favorites).Methods("GET")
	r.HandleFunc("/about", controller.About).Methods("GET")
	r.HandleFunc("/item/{address}/{id}", controller.Id).Methods("GET")

	r.HandleFunc("/login", controller.LoginHandler).Methods("POST")

	r.HandleFunc("/api/favorite", controller.ToggleFavorite).Methods("POST")

	r.HandleFunc("/login/test", controller.LoginTest).Methods("GET")
	r.HandleFunc("/logout", controller.LogoutHandler)

	return r
}
