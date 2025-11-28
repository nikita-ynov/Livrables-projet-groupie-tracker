package main

import (
	"fmt"
	"net/http"
	initTemp "nft/pages"
	"nft/router"
)

func main() {
	initTemp.Init()
	// ... le reste ne change pas
	mux := router.InitRouter()

	port := ":8080"
	fmt.Printf("Serveur lancé sur http://localhost%s\n", port)

	http.ListenAndServe(port, mux)
}
