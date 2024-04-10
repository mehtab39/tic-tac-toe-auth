package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"user-service/db"
)

type GameStats struct {
	Winner string `json:"winner"`
	Loser  string `json:"loser"`
	GameId string `json:"gameId"`
	Draw   bool   `json:"draw"`
}

func UpdateGameStats(w http.ResponseWriter, r *http.Request) {
	var gameStats GameStats
	err := json.NewDecoder(r.Body).Decode(&gameStats)
	if err != nil {
		log.Printf("error retreiving username %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Connect to the database
	dbInstance, err := db.GetInstance()
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}

	// Update player stats
	err = dbInstance.UpdatePlayerStats(gameStats.Winner, gameStats.Loser, gameStats.Draw)

	if err != nil {
		log.Printf("error updating players %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetGameStats(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	username := parts[3]

	// Connect to the database
	dbInstance, err := db.GetInstance()
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}

	// Update player stats
	wins, losses, gamesPlayed, err := dbInstance.GetGameStats(username)

	if err != nil {
		log.Printf("error updating players %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]int{
		"Wins":        wins,
		"Losses":      losses,
		"GamesPlayed": gamesPlayed,
	}

	// Marshal response JSON object
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling JSON: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set Content-Type header and write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		UpdateGameStats(w, r)
	} else if r.Method == http.MethodGet {
		GetGameStats(w, r)
	} else {
		// Return method not allowed for other HTTP methods
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
