package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"poultry-management.com/pkg/api"
	"poultry-management.com/pkg/repo"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbURL := "postgresql://root:secret@localhost:5432/users?sslmode=disable"

	if err := migrateDatabase(dbURL, "up"); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	fmt.Println("Migrations completed successfully.")

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	r := repo.NewRepository(pool)

	router := api.NewAuthHandler(api.Config{JWTSecret: "supposed-to-be-from-env"}, r, gin.Default())
	router.Router.Run(":8080")
}

// For future Refactoring, consider moving the database connection logic to a separate package or module.
func migrateDatabase(dbURL, direction string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}

	migrationsPath := filepath.Join(wd, "internal", "db", "migrations")
	m, err := migrate.New("file:///"+migrationsPath, dbURL)
	if err != nil {
		return fmt.Errorf("creating migrate instance: %w", err)
	}

	switch direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("migrating up: %w", err)
		}
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("migrating down: %w", err)
		}
	case "force":
		version, _, err := m.Version()
		if err != nil {
			return fmt.Errorf("getting current version: %w", err)
		}
		if err := m.Force(int(version)); err != nil {
			return fmt.Errorf("forcing version: %w", err)
		}
	default:
		return fmt.Errorf("invalid migration direction: %s", direction)
	}

	return nil
}
