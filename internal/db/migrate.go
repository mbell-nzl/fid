package db

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	migratesqlite "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func MigrateUp(dbPath, migrationsDir string) error {
	return runMigration(dbPath, migrationsDir, func(m *migrate.Migrate) error {
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("run up migrations: %w", err)
		}
		return nil
	})
}

func MigrateDown(dbPath, migrationsDir string) error {
	return runMigration(dbPath, migrationsDir, func(m *migrate.Migrate) error {
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("run down migrations: %w", err)
		}
		return nil
	})
}

func runMigration(dbPath, migrationsDir string, apply func(*migrate.Migrate) error) error {
	absoluteMigrationsDir, err := filepath.Abs(migrationsDir)
	if err != nil {
		return fmt.Errorf("resolve migrations dir: %w", err)
	}

	conn, err := Open(dbPath)
	if err != nil {
		return err
	}

	driver, err := migratesqlite.WithInstance(conn, &migratesqlite.Config{})
	if err != nil {
		_ = conn.Close()
		return fmt.Errorf("create sqlite migration driver: %w", err)
	}

	migrator, err := migrate.NewWithDatabaseInstance("file://"+absoluteMigrationsDir, driverName, driver)
	if err != nil {
		_ = conn.Close()
		return fmt.Errorf("create migrator: %w", err)
	}

	applyErr := apply(migrator)
	sourceErr, databaseErr := migrator.Close()
	closeErr := errors.Join(sourceErr, databaseErr, conn.Close())
	if applyErr != nil {
		return errors.Join(applyErr, closeErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close migrator: %w", closeErr)
	}

	return nil
}
