package controller

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"nft/database"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/mux"
)

// --- Structures ---

type AlchemyResponse struct {
	NFTs []NFT `json:"nfts"`
}

type NFT struct {
	Contract    Contract `json:"contract"`
	Id          TokenId  `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Media       []Media  `json:"media"`
	Metadata    Metadata `json:"metadata"`
}

type Contract struct {
	Address string `json:"address"`
}

type TokenId struct {
	TokenId string `json:"tokenId"`
}

type Media struct {
	Gateway   string `json:"gateway"`
	Thumbnail string `json:"thumbnail"`
	Format    string `json:"format"`
}

type Metadata struct {
	Attributes []Attribute `json:"attributes"`
}

type Attribute struct {
	TraitType string      `json:"trait_type"`
	Value     interface{} `json:"value"`
}

type PageData struct {
	NFTs        []NFT
	Filters     map[string][]string
	SearchQuery string
	Error       string
	UserAddress string
	CurrentPage string // Pour gérer la classe "selected" dans le menu
}

type LoginRequest struct {
	Address   string `json:"address"`
	Signature string `json:"signature"`
	Nonce     string `json:"nonce"`
}

type FavoritePayload struct {
	ContractAddress string `json:"contract_address"`
	TokenId         string `json:"token_id"`
	ImageUrl        string `json:"image_url"`
	Title           string `json:"title"`
}

const apiKey = "DAOBNtTU5BYuy6upsJXKl"

// --- Fonctions API Alchemy ---

func fetchNFTsForCollection(contractAddress string) ([]NFT, error) {
	url := fmt.Sprintf("https://eth-mainnet.g.alchemy.com/nft/v2/%s/getNFTsForCollection?contractAddress=%s&withMetadata=true", apiKey, contractAddress)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("erreur API Alchemy %d", res.StatusCode)
	}
	body, _ := io.ReadAll(res.Body)
	var response AlchemyResponse
	json.Unmarshal(body, &response)
	return response.NFTs, nil
}

func fetchNFTMetadata(contractAddress, tokenId string) (*NFT, error) {
	url := fmt.Sprintf("https://eth-mainnet.g.alchemy.com/nft/v2/%s/getNFTMetadata?contractAddress=%s&tokenId=%s", apiKey, contractAddress, tokenId)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	var nft NFT
	json.Unmarshal(body, &nft)
	nft.Id.TokenId = tokenId
	nft.Contract.Address = contractAddress
	return &nft, nil
}

func extractFilters(nfts []NFT) map[string][]string {
	tempMap := make(map[string]map[string]bool)
	for _, nft := range nfts {
		for _, attr := range nft.Metadata.Attributes {
			if attr.TraitType == "" || attr.Value == nil {
				continue
			}
			valStr := fmt.Sprintf("%v", attr.Value)
			if tempMap[attr.TraitType] == nil {
				tempMap[attr.TraitType] = make(map[string]bool)
			}
			tempMap[attr.TraitType][valStr] = true
		}
	}
	filters := make(map[string][]string)
	for trait, valuesMap := range tempMap {
		for val := range valuesMap {
			filters[trait] = append(filters[trait], val)
		}
	}
	return filters
}

// --- Handlers & Logique Site ---

func getUserFromSession(r *http.Request) string {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func renderPage(w http.ResponseWriter, filename string, data any) {
	pathsToCheck := []string{
		"pages/" + filename,
		"templates/" + filename,
	}

	var validPath string
	var found bool

	for _, path := range pathsToCheck {
		if _, err := os.Stat(path); err == nil {
			validPath = path
			found = true
			break
		}
	}

	if !found {
		log.Println("❌ ERREUR : Fichier HTML introuvable :", filename)
		http.Error(w, "Fichier introuvable : "+filename, http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles(validPath)
	if err != nil {
		http.Error(w, "Erreur Template : "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

func Home(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "index.html", PageData{
		UserAddress: getUserFromSession(r),
		CurrentPage: "home",
	})
}

func Collections(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	data := PageData{
		SearchQuery: address,
		UserAddress: getUserFromSession(r),
		CurrentPage: "collections",
	}

	if address != "" {
		nfts, err := fetchNFTsForCollection(address)
		if err != nil {
			data.Error = err.Error()
		} else {
			data.NFTs = nfts
			data.Filters = extractFilters(nfts)
		}
	}
	renderPage(w, "collections.html", data)
}

func Id(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userAddr := getUserFromSession(r)

	// 1. Fetch NFT
	nft, err := fetchNFTMetadata(vars["address"], vars["id"])

	// 2. Vérifier si c'est un favori (pour colorer le bouton en rouge)
	isFavorite := false
	if userAddr != "" {
		var count int
		err := database.DB.QueryRow("SELECT COUNT(*) FROM favorites WHERE user_address = ? AND contract_address = ? AND token_id = ?",
			userAddr, vars["address"], vars["id"]).Scan(&count)
		if err == nil && count > 0 {
			isFavorite = true
		}
	}

	data := map[string]interface{}{
		"NFT":         nft,
		"Error":       "",
		"UserAddress": userAddr,
		"IsFavorite":  isFavorite, // Utilisé dans le HTML pour la classe "active"
		"CurrentPage": "",
	}
	if err != nil {
		data["Error"] = err.Error()
	}

	renderPage(w, "id.html", data)
}

// Favorites : Lit depuis la BDD SQL
func Favorites(w http.ResponseWriter, r *http.Request) {
	userAddr := getUserFromSession(r)

	if userAddr == "" {
		renderPage(w, "favorites.html", PageData{
			Error:       "Veuillez connecter votre Wallet",
			CurrentPage: "favorites",
		})
		return
	}

	rows, err := database.DB.Query("SELECT contract_address, token_id, image_url, title FROM favorites WHERE user_address = ?", userAddr)
	if err != nil {
		http.Error(w, "Erreur BDD", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var favNFTs []NFT
	for rows.Next() {
		var f NFT
		var cAddr, tId, img, title string
		rows.Scan(&cAddr, &tId, &img, &title)
		f.Contract.Address = cAddr
		f.Id.TokenId = tId
		f.Title = title
		f.Media = []Media{{Gateway: img}}
		favNFTs = append(favNFTs, f)
	}

	renderPage(w, "favorites.html", PageData{
		NFTs:        favNFTs,
		UserAddress: userAddr,
		CurrentPage: "favorites",
	})
}

func About(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "about.html", PageData{
		UserAddress: getUserFromSession(r),
		CurrentPage: "about",
	})
}

// API : Toggle Favorite (Ajout/Retrait BDD)
func ToggleFavorite(w http.ResponseWriter, r *http.Request) {
	userAddr := getUserFromSession(r)
	if userAddr == "" {
		http.Error(w, "Non connecté", http.StatusUnauthorized)
		return
	}

	var p FavoritePayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Erreur JSON", http.StatusBadRequest)
		return
	}

	_, err := database.DB.Exec("INSERT INTO favorites(user_address, contract_address, token_id, image_url, title) VALUES(?, ?, ?, ?, ?)",
		userAddr, p.ContractAddress, p.TokenId, p.ImageUrl, p.Title)

	if err != nil {
		// Si erreur (déjà existant), on supprime
		database.DB.Exec("DELETE FROM favorites WHERE user_address=? AND contract_address=? AND token_id=?",
			userAddr, p.ContractAddress, p.TokenId)
		w.Write([]byte("removed"))
	} else {
		w.Write([]byte("added"))
	}
}

// LOGIN Handler (Crypto)
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Erreur JSON", http.StatusBadRequest)
		return
	}

	if verifySignature(req.Address, req.Signature, req.Nonce) {
		fmt.Printf("✅ Login OK : %s\n", req.Address)
		database.DB.Exec("INSERT OR IGNORE INTO users (address) VALUES (?)", req.Address)

		expiration := time.Now().Add(24 * time.Hour)
		cookie := &http.Cookie{
			Name:     "session_token",
			Value:    req.Address,
			Path:     "/",
			Expires:  expiration,
			HttpOnly: false,
		}
		http.SetCookie(w, cookie)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Connecté"))
	} else {
		http.Error(w, "Signature invalide", http.StatusUnauthorized)
	}
}

// LOGOUT Handler
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// On écrase le cookie avec une date expirée
	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: false,
	}
	http.SetCookie(w, cookie)

	// On redirige vers l'accueil
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// LoginTest (Pour bouton Discord)
func LoginTest(w http.ResponseWriter, r *http.Request) {
	fakeUser := "Discord_User_Dev"

	// 1. On l'ajoute dans la BDD pour que les favoris marchent
	database.DB.Exec("INSERT OR IGNORE INTO users (address) VALUES (?)", fakeUser)

	// 2. On crée le cookie de session manuellement
	expiration := time.Now().Add(24 * time.Hour)
	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    fakeUser,
		Path:     "/",
		Expires:  expiration,
		HttpOnly: false,
	}
	http.SetCookie(w, cookie)

	// 3. Redirection vers la page Collections
	http.Redirect(w, r, "/collections", http.StatusSeeOther)
}

func verifySignature(walletAddr string, signature string, message string) bool {
	sig, err := hexutil.Decode(signature)
	if err != nil {
		return false
	}
	if len(sig) == 65 && sig[64] >= 27 {
		sig[64] -= 27
	}

	prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	hash := crypto.Keccak256Hash([]byte(prefix))

	sigPublicKey, err := crypto.SigToPub(hash.Bytes(), sig)
	if err != nil {
		return false
	}

	recoveredAddr := crypto.PubkeyToAddress(*sigPublicKey).Hex()
	// Correction ici : strings.EqualFold (plus performant et évite le warning jaune)
	return strings.EqualFold(recoveredAddr, walletAddr)
}
