package main

import (
	"fmt"
	"net/http"
	// Comme ton routeur est dans le dossier src/router, l'import doit être :
	"nft/router" 
)

func main() {
	// ... le reste ne change pas
	mux := router.InitRouter()
	
	port := ":8080"
	fmt.Printf("Serveur lancé sur http://localhost%s\n", port)
	
	http.ListenAndServe(port, mux)
}