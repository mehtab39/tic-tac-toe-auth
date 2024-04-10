package main

import (
	"log"
	"net/http"
	"user-service/db"
	"user-service/handlers"
	"user-service/middlewares"

	"github.com/joho/godotenv"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}
}

func main() {
	_, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Define your routes
	http.HandleFunc("/api/create-account", handlers.SignupHandler)
	http.HandleFunc("/api/signin", handlers.SigninHandler)
	http.HandleFunc("/api/game-stats", handlers.StatsHandler)
	http.HandleFunc("/api/game-stats/", handlers.StatsHandler)
	http.HandleFunc("/ping", pingHandler)

	http.HandleFunc("/api/combine/user/", handlers.GetCombineUserInfo)

	log.Printf("Listening to port %d", 8081)
	http.ListenAndServe(":8081", middlewares.CORSPolicy(http.DefaultServeMux))
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("pong"))
	if err != nil {
		// Handle error if writing to response fails
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
