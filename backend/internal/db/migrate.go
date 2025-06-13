package db

import (
	"database/sql"
	"embed"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	iofs "github.com/golang-migrate/migrate/v4/source/iofs"

	_ "github.com/jackc/pgx/v4/stdlib" // Register pgx driver for database/sql
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// RunMigrations applies all database migrations embedded in the binary.
func RunMigrations() error {
	// Create the source driver from embedded SQL files.
	src, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return fmt.Errorf("creating migration source: %w", err)
	}

	// Open a *sql.DB using pgx stdlib driver.
	dbURL := os.Getenv("DATABASE_URL")
	sqlDB, err := sql.Open("pgx", dbURL)
	if err != nil {
		return fmt.Errorf("opening sql DB: %w", err)
	}

	// Create the database driver for migrate using *sql.DB instance.
	driver, err := pgx.WithInstance(sqlDB, &pgx.Config{})
	if err != nil {
		return fmt.Errorf("creating migration database driver: %w", err)
	}

	// Initialize the migrate instance.
	m, err := migrate.NewWithInstance("iofs", src, "postgres", driver)
	if err != nil {
		return fmt.Errorf("initializing migrate instance: %w", err)
	}

	// Apply all up migrations. Ignore ErrNoChange which means database is up-to-date.
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("applying migrations: %w", err)
	}

	return nil
} 