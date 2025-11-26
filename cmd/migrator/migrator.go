package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/stannisl/url-shortener/internal/infra/postgres"
)

const (
	defaultMigrationsPath = "file://migrations"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Define flags
	var (
		command        = flag.String("command", "", "Migration command: up, down, down-all, version, force, steps")
		migrationsPath = flag.String("path", defaultMigrationsPath, "Path to migrations folder")
		forceVersion   = flag.Int("force-version", 0, "Version to force (use with -command=force)")
		steps          = flag.Int("steps", 0, "Number of steps (use with -command=steps)")
	)

	flag.Parse()

	// Validate command
	if *command == "" {
		printUsage()
		os.Exit(1)
	}

	// Get database URL from environment
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Create migrator
	migrator, err := postgres.NewMigrator(databaseURL, *migrationsPath)
	if err != nil {
		log.Fatalf("Failed to create migrator: %v", err)
	}
	defer migrator.Close()

	// Execute command
	switch *command {
	case "up":
		if err := migrator.Up(); err != nil {
			log.Fatalf("Migration up failed: %v", err)
		}

	case "down":
		if err := migrator.Down(); err != nil {
			log.Fatalf("Migration down failed: %v", err)
		}

	case "down-all":
		if err := migrator.DownAll(); err != nil {
			log.Fatalf("Migration down-all failed: %v", err)
		}

	case "version":
		version, dirty, err := migrator.Version()
		if err != nil {
			log.Fatalf("Failed to get version: %v", err)
		}
		fmt.Printf("Current version: %d, Dirty: %v\n", version, dirty)

	case "force":
		if *forceVersion == 0 {
			log.Fatal("Please specify -force-version flag")
		}
		if err := migrator.Force(*forceVersion); err != nil {
			log.Fatalf("Force version failed: %v", err)
		}

	case "steps":
		if *steps == 0 {
			log.Fatal("Please specify -steps flag (positive for up, negative for down)")
		}
		if err := migrator.Steps(*steps); err != nil {
			log.Fatalf("Steps migration failed: %v", err)
		}

	default:
		printUsage()
		os.Exit(1)
	}

	log.Println("✅ Done!")
}

func printUsage() {
	fmt.Println(`
️Database Migrator

Usage:
  go run cmd/migrator/main.go -command= [options]

Commands:
  up        Apply all pending migrations
  down      Rollback the last migration  
  down-all  Rollback all migrations
  version   Show current migration version
  force     Force set version (use with -force-version)
  steps     Migrate by N steps (use with -steps, negative for rollback)

Options:
  -path           Path to migrations folder (default: file://migrations)
  -force-version  Version number for force command
  -steps          Number of steps for steps command

Examples:
  go run cmd/migrator/main.go -command=up
  go run cmd/migrator/main.go -command=down
  go run cmd/migrator/main.go -command=version
  go run cmd/migrator/main.go -command=force -force-version=1
  go run cmd/migrator/main.go -command=steps -steps=2
  go run cmd/migrator/main.go -command=steps -steps=-1

Environment:
  DATABASE_URL    PostgreSQL connection string (required)
                  Example: postgres://user:pass@localhost:5432/dbname?sslmode=disable`)
}
