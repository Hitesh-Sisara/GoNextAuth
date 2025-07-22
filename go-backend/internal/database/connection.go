package database

import (
	"context"
	"log"

	"github.com/Hitesh-Sisara/GoNextAuth/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB() *pgxpool.Pool {
	config := config.AppConfig

	// Parse database config
	dbConfig, err := pgxpool.ParseConfig(config.Database.URL)
	if err != nil {
		log.Fatalf("Failed to parse database config: %v", err)
	}

	// Set connection pool settings
	dbConfig.MaxConns = 30
	dbConfig.MinConns = 5

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		log.Fatalf("Failed to create connection pool: %v", err)
	}

	// Test connection
	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to database")

	DB = pool
	return pool
}

func GetDB() *pgxpool.Pool {
	return DB
}

func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed")
	}
}

// Helper function for single connection operations
func GetConn(ctx context.Context) (*pgxpool.Conn, error) {
	return DB.Acquire(ctx)
}
