CREATE TABLE player_stats (
    playerID VARCHAR(255) REFERENCES users(username),
    Wins INT DEFAULT 0,
    Loss INT DEFAULT 0,
    GamesPlayed INT DEFAULT 0,
    PRIMARY KEY (playerID)
);



GRANT ALL PRIVILEGES ON TABLE player_stats TO mehtabgill;

