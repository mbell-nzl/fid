package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

const driverName = "sqlite"

var requiredPragmas = []string{
	"PRAGMA foreign_keys = ON;",
	"PRAGMA recursive_triggers = ON;",
}

func Open(path string) (*sql.DB, error) {
	conn, err := sql.Open(driverName, path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}

	for _, pragma := range requiredPragmas {
		if _, err := conn.Exec(pragma); err != nil {
			_ = conn.Close()
			return nil, fmt.Errorf("apply pragma %q: %w", pragma, err)
		}
	}

	if err := conn.Ping(); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("ping sqlite database: %w", err)
	}

	return conn, nil
}
