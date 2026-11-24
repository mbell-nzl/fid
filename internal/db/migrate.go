package db

import (
	"errors"
	"fmt"
	"io/fs"

	"github.com/golang-migrate/migrate/v4"
	migratesqlite "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func MigrateUp(dbPath string, migrationsFS fs.FS) error {
	return runMigration(dbPath, migrationsFS, func(m *migrate.Migrate) error {
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("run up migrations: %w", err)
		}
		return nil
	})
}

func MigrateDown(dbPath string, migrationsFS fs.FS) error {
	return runMigration(dbPath, migrationsFS, func(m *migrate.Migrate) error {
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("run down migrations: %w", err)
		}
		return nil
	})
}

func runMigration(dbPath string, migrationsFS fs.FS, apply func(*migrate.Migrate) error) error {
	if migrationsFS == nil {
		return fmt.Errorf("migrations filesystem is required")
	}
	source, err := iofs.New(migrationsFS, ".")
	if err != nil {
		return fmt.Errorf("open migration source: %w", err)
	}

	conn, err := Open(dbPath)
	if err != nil {
		_ = source.Close()
		return err
	}

	driver, err := migratesqlite.WithInstance(conn, &migratesqlite.Config{})
	if err != nil {
		_ = source.Close()
		_ = conn.Close()
		return fmt.Errorf("create sqlite migration driver: %w", err)
	}

	migrator, err := migrate.NewWithInstance("iofs", source, driverName, driver)
	if err != nil {
		_ = source.Close()
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
