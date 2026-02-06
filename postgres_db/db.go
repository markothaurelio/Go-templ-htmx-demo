package postgres_db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // Import the pgx driver
	"github.com/rs/zerolog/log"
)

// Database configuration (PLACEHOLDER VALUES)
// Configure locally or via environment-specific files
const (
	DB_USER    = "DB_USER"
	DB_PASS    = "DB_PASSWORD"
	DB_NAME    = "DB_NAME"
	DB_HOST    = "DB_HOST"
	DB_PORT    = "DB_PORT"
	DB_SSLMODE = "DB_SSLMODE"
)

var DB *sql.DB

// Connect to PostgreSQL
func ConnectDB() *sql.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		DB_USER, DB_PASS, DB_HOST, DB_PORT, DB_NAME, DB_SSLMODE)

	var err error
	DB, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open database")
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := DB.PingContext(ctx); err != nil {
		log.Fatal().Err(err).Msg("Database ping failed")
	}

	log.Info().Msg("✅ Connected to PostgreSQL database")
	return DB
}

// Close the database connection
func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Info()
	}
}
