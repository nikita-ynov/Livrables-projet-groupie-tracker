package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // Driver MySQL
)

var DB *sql.DB

func InitDB() {
	// --- CONFIGURATION ALWAYS DATA ---
	// Remplace les valeurs ci-dessous par celles de ton tableau de bord AlwaysData > MySQL

	username := "443067"                          // Ex: tonpseudo
	password := "giogio220706"                    // Le mot de passe que tu as défini
	host := "mysql-borderlandsapi.alwaysdata.net" // L'adresse de l'hôte
	dbName := "borderlandsapi_nft"                // Le nom exact de ta base sur AlwaysData

	// Construction de la chaîne de connexion
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", username, password, host, dbName)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Erreur de configuration (DSN) :", err)
	}

	// Test de la connexion vers Internet
	if err = DB.Ping(); err != nil {
		log.Fatal("Impossible de joindre AlwaysData (Vérifie l'accès distant et le mot de passe !) :", err)
	}

	fmt.Println("✅ Connecté avec succès à la base de données AlwaysData !")

	createTables()
}

func createTables() {
	// 1. Table USERS
	createUsers := `
    CREATE TABLE IF NOT EXISTS users (
        id INT AUTO_INCREMENT PRIMARY KEY,
        address VARCHAR(255) NOT NULL UNIQUE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

	// 2. Table FAVORITES
	createFavs := `
    CREATE TABLE IF NOT EXISTS favorites (
        id INT AUTO_INCREMENT PRIMARY KEY,
        user_address VARCHAR(255) NOT NULL,
        contract_address VARCHAR(255) NOT NULL,
        token_id VARCHAR(255) NOT NULL,
        image_url TEXT,
        title TEXT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        
        -- Empêche les doublons pour un même utilisateur
        UNIQUE KEY unique_fav (user_address, contract_address, token_id)
    );`

	if _, err := DB.Exec(createUsers); err != nil {
		log.Fatal("Erreur création table users :", err)
	}
	if _, err := DB.Exec(createFavs); err != nil {
		log.Fatal("Erreur création table favorites :", err)
	}

	fmt.Println("✅ Tables 'users' et 'favorites' vérifiées sur le Cloud.")
}
