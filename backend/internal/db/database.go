package db

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	DB  *pgxpool.Pool
	RDB *redis.Client
	Ctx = context.Background()
)

// Connect initializes the database and Redis connections.
func Connect() error {
	// Connect to PostgreSQL
	dbpool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}
	DB = dbpool

	// Connect to Redis
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	RDB = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Ping Redis to check the connection
	if _, err := RDB.Ping(Ctx).Result(); err != nil {
		return fmt.Errorf("unable to connect to redis: %v", err)
	}

	fmt.Println("Successfully connected to PostgreSQL and Redis.")
	return nil
}

// CreateLinkTable creates the 'links' table if it doesn't exist.
func CreateLinkTable() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS links (
		id SERIAL PRIMARY KEY,
		short_code VARCHAR(255) UNIQUE NOT NULL,
		long_url TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		expires_at TIMESTAMP WITH TIME ZONE
	);`

	_, err := DB.Exec(Ctx, createTableSQL)
	if err != nil {
		return fmt.Errorf("error creating links table: %v", err)
	}

	fmt.Println("Links table created or already exists.")
	return nil
} 