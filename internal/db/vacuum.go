package db

import (
	"database/sql"
	"fmt"
)

// Vacuum compacts the SQLite database connected by conn.
func Vacuum(conn *sql.DB) error {
	if conn == nil {
		return fmt.Errorf("vacuum sqlite database: connection is required")
	}
	if _, err := conn.Exec(`VACUUM`); err != nil {
		return fmt.Errorf("vacuum sqlite database: %w", err)
	}
	return nil
}
