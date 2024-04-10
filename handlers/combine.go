package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"user-service/db"
)

func GetCombineUserInfo(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	username := parts[4]

	log.Printf("retrieve user info for %s", username)

	// Get database instance
	database, err := db.GetInstance()
	if err != nil {
		http.Error(w, "Failed to get database instance", http.StatusInternalServerError)
		return
	}

	userInfo, err := database.GetUser(username)
	if err != nil {
		http.Error(w, "Failed to retrieve user information", http.StatusInternalServerError)
		return
	}

	wins, losses, gamesPlayed, err := database.GetGameStats(username)

	if err != nil {
		log.Printf("error updating players %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stats := map[string]int{
		"Wins":        wins,
		"Losses":      losses,
		"GamesPlayed": gamesPlayed,
	}
	response := map[string]interface{}{
		"stats": stats,
		"auth":  userInfo,
	}

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
