// db/database.go

package db

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

// User represents a user in the system
type User struct {
	ID       int
	Username string
}

// Database represents the database connection
type Database struct {
	connection *sql.DB
}

var (
	instance *Database
	once     sync.Once
)

// Connect initializes the database connection if it hasn't been initialized yet
func Connect() (*Database, error) {
	once.Do(func() {
		// Read environment variables
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbName := os.Getenv("DB_NAME")
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")

		// Construct database connection string
		dbURI := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
			dbHost, dbPort, dbName, dbUser, dbPassword)

		// Connect to the database
		db, err := sql.Open("postgres", dbURI)
		if err != nil {
			fmt.Println("error connecting to the database:", err)
			return
		}

		// Initialize the singleton instance
		instance = &Database{
			connection: db,
		}
	})

	return instance, nil
}

// GetInstance returns the singleton instance of the database connection
func GetInstance() (*Database, error) {
	if instance == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}
	return instance, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.connection.Close()
}

// AuthenticateUser checks if the provided username and password are valid
func (d *Database) AuthenticateUser(username, password string) (bool, error) {
	// Prepare the SQL statement to check if the user exists
	query := `
        SELECT COUNT(*) FROM users WHERE username=$1 AND password=$2
    `

	var count int
	err := d.connection.QueryRow(query, username, password).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error authenticating user: %v", err)
	}

	// If count is greater than zero, the user exists and the password is correct
	return count > 0, nil
}

// SaveUser saves a new user to the database
func (d *Database) SaveUser(username, password string) error {
	// Prepare the SQL statement to insert a new user
	query := `
	   INSERT INTO users (username, password) VALUES ($1, $2)
    `

	_, err := d.connection.Exec(query, username, password)
	if err != nil {
		return err
	}

	return nil
}

// getUser retrieves user information from the database based on the username
func (d *Database) GetUser(username string) (*User, error) {
	// Prepare the SQL statement to retrieve user information
	query := `
        SELECT id, username FROM users WHERE username=$1
    `

	// Execute the SQL query
	row := d.connection.QueryRow(query, username)

	// Initialize variables to store the retrieved user information
	var user User

	// Scan the row into the user struct
	err := row.Scan(&user.ID, &user.Username)
	if err != nil {
		// If the user doesn't exist or an error occurs, return an appropriate error
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error retrieving user: %v", err)
	}

	// Return the retrieved user
	return &user, nil
}
