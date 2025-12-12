package main

import (
	"fmt"
	"net/http"
	"nft/database"
	"nft/router"
)

func main() {
	// 1. Démarrer la base de données
	database.InitDB()

	// 2. Configurer les routes
	r := router.InitRouter()

	fmt.Println("🚀 Serveur lancé sur http://localhost:8080")

	// 3. Lancer le serveur
	http.ListenAndServe(":8080", r)
}
