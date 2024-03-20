package handlers

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/davidroman0O/gogog/data"
	"github.com/davidroman0O/gogog/types"
)

type GameTable struct {
	ID         int      `json:"id"`
	Title      string   `json:"title"`
	Category   string   `json:"category"`
	Platforms  []string `json:"platforms"`
	Downloaded bool     `json:"downloaded"`
}

func filterAndConvertProducts(products []types.Product, categoryFilter string) []GameTable {
	var games []GameTable

	for _, prod := range products {
		if prod.IsGame && strings.EqualFold(prod.Category, categoryFilter) {
			var platforms []string
			if prod.WorksOn.Windows {
				platforms = append(platforms, "Windows")
			}
			if prod.WorksOn.Mac {
				platforms = append(platforms, "Mac")
			}
			if prod.WorksOn.Linux {
				platforms = append(platforms, "Linux")
			}

			game := GameTable{
				ID:         prod.ID,
				Title:      prod.Title,
				Category:   prod.Category,
				Platforms:  platforms,
				Downloaded: false, // Assuming true for the example
			}
			games = append(games, game)
		}
	}

	return games
}

func GetGamesTablefunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", "application/json")
	products, err := data.GetGames()
	if err != nil {
		http.Redirect(w, r, "/404", http.StatusFound)
		return
	}

	draw := r.URL.Query().Get("draw")
	search := r.URL.Query().Get("search[value]")
	start, _ := strconv.Atoi(r.URL.Query().Get("start"))
	length, _ := strconv.Atoi(r.URL.Query().Get("length"))
	// Extract DataTables ordering parameters
	orderColumnIndex, _ := strconv.Atoi(r.URL.Query().Get("order[0][column]"))
	orderDir := r.URL.Query().Get("order[0][dir]")

	// Filter and convert to GameTable based on search criteria
	filteredGames := filterAndConvertProducts(products, search)

	// Apply sorting to filteredGames based on DataTables request
	sort.Slice(filteredGames, func(i, j int) bool {
		switch orderColumnIndex {
		case 0: // Title
			if orderDir == "asc" {
				return strings.ToLower(filteredGames[i].Title) < strings.ToLower(filteredGames[j].Title)
			}
			return strings.ToLower(filteredGames[i].Title) > strings.ToLower(filteredGames[j].Title)
		case 1: // Category
			if orderDir == "asc" {
				return filteredGames[i].Category < filteredGames[j].Category
			}
			return filteredGames[i].Category > filteredGames[j].Category
		default:
			return true // Default no sorting if column index is out of range
		}
	})

	// Implement pagination
	var paginatedGames []GameTable
	if start+length > len(filteredGames) {
		paginatedGames = filteredGames[start:]
	} else {
		paginatedGames = filteredGames[start : start+length]
	}

	// Prepare the response for DataTables
	response := map[string]interface{}{
		"draw":            draw,
		"recordsTotal":    len(products),      // Total number of products before filtering
		"recordsFiltered": len(filteredGames), // Total number of products after filtering
		"data":            paginatedGames,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
