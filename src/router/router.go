package router

import (
	"net/http"
	"nft/controller" 
)

func InitRouter() *http.ServeMux {
	mux := http.NewServeMux()
	
	fileServer := http.FileServer(http.Dir("assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fileServer))

	// Route Index
	mux.HandleFunc("/", controller.IndexHandler)

	return mux
}