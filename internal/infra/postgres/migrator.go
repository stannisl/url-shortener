package postgres

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Migrator handles database migrations
type Migrator struct {
	migrate *migrate.Migrate
}

// NewMigrator creates a new migrator instance
func NewMigrator(databaseURL, migrationsPath string) (*Migrator, error) {
	if databaseURL == "" {
		return nil, errors.New("database URL is required")
	}

	if migrationsPath == "" {
		migrationsPath = "file://migrations"
	}

	m, err := migrate.New(migrationsPath, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator: %w", err)
	}

	return &Migrator{migrate: m}, nil
}

// Up applies all available migrations
func (m *Migrator) Up() error {
	log.Println("Applying migrations...")

	err := m.migrate.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No new migrations to apply")
			return nil
		}
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Println("Migrations applied successfully")
	return nil
}

// Down rolls back the last migration
func (m *Migrator) Down() error {
	log.Println("Rolling back last migration...")

	err := m.migrate.Steps(-1)
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No migrations to rollback")
			return nil
		}
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	log.Println("Migration rolled back successfully")
	return nil
}

// DownAll rolls back all migrations
func (m *Migrator) DownAll() error {
	log.Println("Rolling back all migrations...")

	err := m.migrate.Down()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No migrations to rollback")
			return nil
		}
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	log.Println("All migrations rolled back successfully")
	return nil
}

// Version returns the current migration version
func (m *Migrator) Version() (uint, bool, error) {
	version, dirty, err := m.migrate.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			return 0, false, nil // no migrations applied yet
		}
		return 0, false, fmt.Errorf("failed to get version: %w", err)
	}
	return version, dirty, nil
}

// Force sets the migration version without running migrations
// Useful for fixing dirty state
func (m *Migrator) Force(version int) error {
	log.Printf("Forcing version to %d...\n", version)

	err := m.migrate.Force(version)
	if err != nil {
		return fmt.Errorf("failed to force version: %w", err)
	}

	log.Println("Version forced successfully")
	return nil
}

// Steps migrates up or down by the specified number of steps
func (m *Migrator) Steps(n int) error {
	if n > 0 {
		log.Printf("Applying %d migration(s)...\n", n)
	} else {
		log.Printf("Rolling back %d migration(s)...\n", -n)
	}

	err := m.migrate.Steps(n)
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No migrations to apply")
			return nil
		}
		return fmt.Errorf("failed to apply steps: %w", err)
	}

	log.Println("Steps applied successfully")
	return nil
}

// Close closes the migrator
func (m *Migrator) Close() error {
	sourceErr, dbErr := m.migrate.Close()
	if sourceErr != nil {
		return fmt.Errorf("source close error: %w", sourceErr)
	}
	if dbErr != nil {
		return fmt.Errorf("database close error: %w", dbErr)
	}
	return nil
}
