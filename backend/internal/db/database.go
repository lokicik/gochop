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

// CreateUsersTable creates the NextAuth.js compatible 'users' table if it doesn't exist.
func CreateUsersTable() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR(255),
		email VARCHAR(255) UNIQUE NOT NULL,
		email_verified TIMESTAMPTZ,
		image TEXT,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		is_admin BOOLEAN DEFAULT FALSE
	);`

	_, err := DB.Exec(Ctx, createTableSQL)
	if err != nil {
		return fmt.Errorf("error creating users table: %v", err)
	}

	fmt.Println("Users table created or already exists.")
	return nil
}

// CreateAccountsTable creates the NextAuth.js 'accounts' table for OAuth providers.
func CreateAccountsTable() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS accounts (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		type VARCHAR(255) NOT NULL,
		provider VARCHAR(255) NOT NULL,
		provider_account_id VARCHAR(255) NOT NULL,
		refresh_token TEXT,
		access_token TEXT,
		expires_at BIGINT,
		token_type VARCHAR(255),
		scope TEXT,
		id_token TEXT,
		session_state TEXT,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(provider, provider_account_id)
	);`

	_, err := DB.Exec(Ctx, createTableSQL)
	if err != nil {
		return fmt.Errorf("error creating accounts table: %v", err)
	}

	fmt.Println("Accounts table created or already exists.")
	return nil
}

// CreateSessionsTable creates the NextAuth.js 'sessions' table for session management.
func CreateSessionsTable() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS sessions (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		session_token VARCHAR(255) UNIQUE NOT NULL,
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		expires TIMESTAMPTZ NOT NULL,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := DB.Exec(Ctx, createTableSQL)
	if err != nil {
		return fmt.Errorf("error creating sessions table: %v", err)
	}

	fmt.Println("Sessions table created or already exists.")
	return nil
}

// CreateVerificationTokensTable creates the NextAuth.js 'verification_tokens' table for email verification.
func CreateVerificationTokensTable() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS verification_tokens (
		identifier VARCHAR(255) NOT NULL,
		token VARCHAR(255) UNIQUE NOT NULL,
		expires TIMESTAMPTZ NOT NULL,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY(identifier, token)
	);`

	_, err := DB.Exec(Ctx, createTableSQL)
	if err != nil {
		return fmt.Errorf("error creating verification_tokens table: %v", err)
	}

	fmt.Println("Verification tokens table created or already exists.")
	return nil
}

// CreateLinkTable creates the 'links' table if it doesn't exist.
func CreateLinkTable() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS links (
		id SERIAL PRIMARY KEY,
		short_code VARCHAR(255) UNIQUE NOT NULL,
		long_url TEXT NOT NULL,
		context TEXT,
		user_id UUID REFERENCES users(id),
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

// CreateAnalyticsTable creates the 'analytics' table if it doesn't exist.
func CreateAnalyticsTable() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS analytics (
		id SERIAL PRIMARY KEY,
		short_code VARCHAR(255) NOT NULL,
		ip_address INET,
		user_agent TEXT,
		referrer TEXT,
		country VARCHAR(255),
		region VARCHAR(255),
		city VARCHAR(255),
		clicked_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (short_code) REFERENCES links(short_code) ON DELETE CASCADE
	);`

	_, err := DB.Exec(Ctx, createTableSQL)
	if err != nil {
		return fmt.Errorf("error creating analytics table: %v", err)
	}

	fmt.Println("Analytics table created or already exists.")
	return nil
} 