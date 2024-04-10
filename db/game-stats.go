package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type GameStats struct {
	Username string
	Status   string
}

func (d *Database) UpdatePlayerStats(winner string, loser string, draw bool) error {
	// Define the SQL statements
	sqlStmt := map[string]string{
		"winner": `
			INSERT INTO player_stats (playerID, Wins, Loss, GamesPlayed)
			VALUES ($1, 1, 0, 1)
			ON CONFLICT (playerID) DO UPDATE
			SET Wins = player_stats.Wins + 1, GamesPlayed = player_stats.GamesPlayed + 1
			WHERE player_stats.playerID = $1;
		`,
		"loser": `
			INSERT INTO player_stats (playerID, Wins, Loss, GamesPlayed)
			VALUES ($1, 0, 1, 1)
			ON CONFLICT (playerID) DO UPDATE
			SET Loss = player_stats.Loss + 1, GamesPlayed = player_stats.GamesPlayed + 1
			WHERE player_stats.playerID = $1;
		`,
	}

	// Prepare the SQL statements
	stmt := make(map[string]*sql.Stmt)
	for key, value := range sqlStmt {
		statement, err := d.connection.Prepare(value)
		if err != nil {
			return err
		}
		defer statement.Close()
		stmt[key] = statement
	}

	// Execute the appropriate SQL statement for winner
	_, err := stmt["winner"].Exec(winner)
	if err != nil {
		return err
	}

	// Execute the appropriate SQL statement for loser
	_, err = stmt["loser"].Exec(loser)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) GetGameStats(username string) (int, int, int, error) {
	// Define the SQL query
	sqlQuery := `
        SELECT Wins, Loss, GamesPlayed
        FROM player_stats
        WHERE playerID = $1
    `

	// Execute the SQL query
	row := d.connection.QueryRow(sqlQuery, username)

	// Variables to store the retrieved values
	var wins, loss, gamesPlayed int

	// Scan the values from the row into variables
	err := row.Scan(&wins, &loss, &gamesPlayed)
	if err != nil {
		if err == sql.ErrNoRows {
			// If no rows are returned, it means the user doesn't exist in player_stats
			return 0, 0, 0, fmt.Errorf("no game stats found for user %s", username)
		}
		// Return the error if it's not ErrNoRows
		return 0, 0, 0, err
	}

	// Return the retrieved values
	return wins, loss, gamesPlayed, nil
}
