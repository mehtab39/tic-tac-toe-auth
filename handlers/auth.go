// handlers/auth.go

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"user-service/db"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get database instance
	database, err := db.GetInstance()
	if err != nil {
		http.Error(w, "Failed to get database instance", http.StatusInternalServerError)
		return
	}

	// Save user to database
	err = database.SaveUser(user.Username, user.Password)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			http.Error(w, "Username already in use", http.StatusForbidden)
			return
		}
		http.Error(w, "Failed to save user to database", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User signed up successfully")
}

func SigninHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get database instance
	database, err := db.GetInstance()
	if err != nil {
		http.Error(w, "Failed to get database instance", http.StatusInternalServerError)
		return
	}

	// Authenticate user
	authenticated, err := database.AuthenticateUser(user.Username, user.Password)
	if err != nil {
		http.Error(w, "Failed to authenticate user", http.StatusInternalServerError)
		return
	}

	if authenticated {
		// Respond with success message
		userInfo, err := database.GetUser(user.Username)
		if err != nil {
			http.Error(w, "Failed to retrieve user information", http.StatusInternalServerError)
			return
		}

		// Marshal userInfo into JSON format
		userInfoJSON, err := json.Marshal(userInfo)
		if err != nil {
			http.Error(w, "Failed to marshal user information to JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)

		w.Write(userInfoJSON)
	} else {
		// Respond with error message
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
	}
}
