package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"nft/pages"

	"github.com/gorilla/mux"
)

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
}

const apiKey = "DAOBNtTU5BYuy6upsJXKl"

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
		return nil, fmt.Errorf("erreur API Alchemy (Code %d). Vérifie ta clé API ou l'adresse du contrat", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)

	var response AlchemyResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.NFTs, nil
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

// --- Handlers ---

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
	address := r.URL.Query().Get("address")

	data := PageData{
		SearchQuery: address,
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
	id := vars["id"]
	fmt.Println("Item ID = " + id)
	renderPage(w, "id.html", nil)
}

func Favorites(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "favorites.html", nil)
}

func About(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "about.html", nil)
}
