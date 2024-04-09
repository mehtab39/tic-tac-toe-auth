// db/database.go

package db

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

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
